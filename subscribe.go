package mqtt

import (
	"context"
	"encoding/json"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	message paho.Message
	vars    []string
}

type MessageHandler func(Message)

func (m *Message) TopicVars() []string {
	return m.vars
}

func (m *Message) Topic() string {
	return m.message.Topic()
}

func (m *Message) QOS() QOS {
	return QOS(m.message.Qos())
}

func (m *Message) IsDuplicate() bool {
	return m.message.Duplicate()
}

func (m *Message) Acknowledge() {
	m.message.Ack()
}

func (m *Message) Payload() []byte {
	return m.message.Payload()
}

func (m *Message) PayloadString() string {
	return string(m.message.Payload())
}

func (m *Message) PayloadJSON(v interface{}) error {
	return json.Unmarshal(m.message.Payload(), v)
}

func (c *Client) Handle(topic string, handler MessageHandler) Route {
	return c.router.addRoute(topic, handler)
}

func (c *Client) Listen(topic string) (chan Message, Route) {
	queue := make(chan Message)
	route := c.router.addRoute(topic, func(message Message) {
		queue <- message
	})
	return queue, route
}

func (c *Client) Subscribe(ctx context.Context, topic string, qos QOS) error {
	token := c.client.Subscribe(topic, byte(qos), nil)
	err := tokenWithContext(ctx, token)
	return err
}

func (c *Client) Unsubscribe(ctx context.Context, topic string) error {
	token := c.client.Unsubscribe(topic)
	err := tokenWithContext(ctx, token)
	return err
}
