package mqtt

import (
	"context"
	"encoding/json"
)

// PublishOption are extra options when publishing a message
type PublishOption int

const (
	// Retain tells the broker to retain a message and send it as the first message to new subscribers.
	Retain PublishOption = iota
)

// Publish a message with a byte array payload
func (c *Client) Publish(ctx context.Context, topic string, payload []byte, qos QOS, options ...PublishOption) error {
	return c.publish(ctx, topic, payload, qos, options)
}

// PublishString publishes a message with a string payload
func (c *Client) PublishString(ctx context.Context, topic string, payload string, qos QOS, options ...PublishOption) error {
	return c.publish(ctx, topic, []byte(payload), qos, options)
}

// PublishJSON publishes a message with the payload encoded as JSON using encoding/json
func (c *Client) PublishJSON(ctx context.Context, topic string, payload interface{}, qos QOS, options ...PublishOption) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.publish(ctx, topic, data, qos, options)
}

func (c *Client) publish(ctx context.Context, topic string, payload []byte, qos QOS, options []PublishOption) error {
	var retained = false
	for _, option := range options {
		switch option {
		case Retain:
			retained = true
		}
	}

	token := c.client.Publish(topic, byte(qos), retained, payload)
	return tokenWithContext(ctx, token)
}
