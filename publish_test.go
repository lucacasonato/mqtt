package mqtt_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucacasonato/mqtt"
)

var testUUID = uuid.New().String()

// TestPublishSuccess checks that a message publish succeeds
func TestPublishSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	err = client.Publish(ctx(), testUUID+"/TestPublishSuccess", []byte("hello"), mqtt.AtLeastOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
}

// TestPublishContextTimeout checks that a message publish errors if a context with a timeout times out
func TestPublishContextTimeout(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx(), 1*time.Nanosecond)
	defer cancel()
	err = client.Publish(ctx, testUUID+"/TestPublishContextTimeout", []byte("hello"), mqtt.AtLeastOnce)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("publish should have returned the error context.DeadlineExceeded")
	}
}

// TestPublishContextCancelled checks that a message publish errors if a context with a cancel gets canceled
func TestPublishContextCancelled(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	ctx, cancel := context.WithCancel(ctx())
	go func() {
		time.Sleep(1 * time.Microsecond)
		cancel()
	}()
	defer cancel()
	err = client.Publish(ctx, testUUID+"/TestPublishContextCancelled", []byte("hello"), mqtt.AtLeastOnce)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("publish should have returned the error context.Canceled")
	}
}

// TestPublishFailed checks that a invalid publish does not get publish but errors
func TestPublishFailed(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	err = client.Publish(ctx(), testUUID+"/TestPublishFailed", nil, 3)
	if err == nil {
		t.Fatalf("publish should have failed")
	}
}

// TestPublishSuccess checks that a message publish succeeds
func TestPublishSuccessRetained(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	err = client.Publish(ctx(), testUUID+"/TestPublishSuccessRetained", []byte("hello"), mqtt.AtLeastOnce, mqtt.Retain)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
}

// TestPublisStringSuccess checks that a string message publish succeeds
func TestPublisStringSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	err = client.PublishString(ctx(), testUUID+"/TestPublisStringSuccess", "world", mqtt.AtLeastOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
}

// TestPublisJSONSuccess checks that json message publish succeeds
func TestPublisJSONSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	err = client.PublishJSON(ctx(), testUUID+"/TestPublisJSONSuccess", []string{"hello", "world"}, mqtt.AtLeastOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
}

// TestPublisJSONFailed checks that json message fails to parse
func TestPublisJSONFailed(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(ctx())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	err = client.PublishJSON(ctx(), testUUID+"/TestPublisJSONFailed", make(chan int), mqtt.AtLeastOnce)
	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Fatalf("publish error should be of type *json.UnsupportedTypeError: %v", err)
	}
}
