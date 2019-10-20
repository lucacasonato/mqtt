package mqtt

import (
	"context"
	"encoding/json"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	message paho.Message
}

type MessageHandler func(Message)

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

func (c *Client) Listen(handler MessageHandler, topics ...string) {
	for _, topic := range topics {
		c.router.addRoute(topic, handler)
	}
}

func (c *Client) Subscribe(ctx context.Context, handler MessageHandler, topic string, qos QOS) error {
	token := c.client.Subscribe(topic, byte(qos), nil)
	err := tokenWithContext(ctx, token)
	if err != nil {
		return err
	}
	c.router.addRoute(topic, handler)
	return nil
}
