package mqtt

import (
	"context"
	"encoding/json"
)

type PublishOption int

const (
	Retain PublishOption = iota
)

func (c *Client) Publish(ctx context.Context, topic string, payload []byte, qos QOS, options ...PublishOption) error {
	return c.publish(ctx, topic, payload, qos, options)
}

func (c *Client) PublishString(ctx context.Context, topic string, payload string, qos QOS, options ...PublishOption) error {
	return c.publish(ctx, topic, []byte(payload), qos, options)
}

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
