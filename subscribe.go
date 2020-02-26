package mqtt

import (
	"context"
	"encoding/json"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// A Message from or to the broker
type Message struct {
	message paho.Message
	vars    []string
}

// A MessageHandler to handle incoming messages
type MessageHandler func(Message)

// TopicVars is a list of all the message specific matches for a wildcard in a route topic.
// If the route would be `config/+/full` and the messages topic is `config/server_1/full` then thous would return `[]string{"server_1"}`
func (m *Message) TopicVars() []string {
	return m.vars
}

// Topic is the topic the message was recieved on
func (m *Message) Topic() string {
	return m.message.Topic()
}

// QOS is the quality of service the message was recieved with
func (m *Message) QOS() QOS {
	return QOS(m.message.Qos())
}

// IsDuplicate is true if this exact message has been recieved before (due to a AtLeastOnce QOS)
func (m *Message) IsDuplicate() bool {
	return m.message.Duplicate()
}

// Acknowledge explicitly acknowledges to a broker that the message has been recieved
func (m *Message) Acknowledge() {
	m.message.Ack()
}

// Payload returns the payload as a byte array
func (m *Message) Payload() []byte {
	return m.message.Payload()
}

// PayloadString returns the payload as a string
func (m *Message) PayloadString() string {
	return string(m.message.Payload())
}

// PayloadJSON unmarshals the payload into the provided interface using encoding/json and returns an error if anything fails
func (m *Message) PayloadJSON(v interface{}) error {
	return json.Unmarshal(m.message.Payload(), v)
}

// Handle adds a handler for a certain topic. This handler gets called if any message arrives that matches the topic.
// Also returns a route that can be used to unsubsribe. Does not automatically subscribe.
func (c *Client) Handle(topic string, handler MessageHandler) Route {
	return c.router.addRoute(topic, handler)
}

// Listen returns a stream of messages that match the topic.
// Also returns a route that can be used to unsubsribe. Does not automatically subscribe.
func (c *Client) Listen(topic string) (chan Message, Route) {
	queue := make(chan Message)
	route := c.router.addRoute(topic, func(message Message) {
		queue <- message
	})
	return queue, route
}

// Subscribe subscribes to a certain topic and errors if this fails.
func (c *Client) Subscribe(ctx context.Context, topic string, qos QOS) error {
	token := c.client.Subscribe(topic, byte(qos), nil)
	err := tokenWithContext(ctx, token)
	return err
}

// SubscribeMultiple subscribes to multiple topics and errors if this fails.
func (c *Client) SubscribeMultiple(ctx context.Context, subscriptions map[string]QOS) error {
	subs := make(map[string]byte, len(subscriptions))
	for topic, qos := range subscriptions {
		subs[topic] = byte(qos)
	}
	token := c.client.SubscribeMultiple(subs, nil)
	err := tokenWithContext(ctx, token)
	return err
}

// Unsubscribe unsubscribes from a certain topic and errors if this fails.
func (c *Client) Unsubscribe(ctx context.Context, topic string) error {
	token := c.client.Unsubscribe(topic)
	err := tokenWithContext(ctx, token)
	return err
}
