package mqtt_test

import (
	"context"
	"testing"
	"time"

	"github.com/lucacasonato/mqtt"
)

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(context.Background())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}

	reciever := make(chan mqtt.Message)
	err = client.Subscribe(context.Background(), func(message mqtt.Message) {
		reciever <- message
	}, testUUID+"/TestSubcribeSuccess", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(context.Background(), testUUID+"/TestSubcribeSuccess", "[1, 2]", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-reciever
	if string(message.Payload()) != "[1, 2]" {
		t.Fatalf("message payload should have been byte array '%v' but is %v", []byte("[1, 2]"), message.Payload())
	}
	if message.PayloadString() != "[1, 2]" {
		t.Fatalf("message payload should have been '[1, 2]' but is %v", message.PayloadString())
	}
	v := []int{}
	err = message.PayloadJSON(&v)
	if err != nil {
		t.Fatalf("json should have unmarshalled: %v", err)
	}
	if len(v) != 2 || v[0] != 1 || v[1] != 2 {
		t.Fatalf("message payload should have been []int{1, 2} but is %v", v)
	}
	if message.Topic() != testUUID+"/TestSubcribeSuccess" {
		t.Fatalf("message topic should be %v but is %v", testUUID+"/TestSubcribeSuccess", message.Topic())
	}
	if message.QOS() != mqtt.ExactlyOnce {
		t.Fatalf("message qos should be mqtt.ExactlyOnce but is %v", message.QOS())
	}
	if message.IsDuplicate() != false {
		t.Fatalf("message IsDuplicate should be false but is %v", message.IsDuplicate())
	}
	message.Acknowledge()
}

// TestListenSuccess checks that a listener recieves a message correctly
func TestListenSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(context.Background())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	reciever := make(chan mqtt.Message)
	err = client.Subscribe(context.Background(), func(message mqtt.Message) {}, testUUID+"/TestListenSuccess", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	client.Listen(func(message mqtt.Message) {
		reciever <- message
	}, testUUID+"/TestListenSuccess")
	err = client.PublishString(context.Background(), testUUID+"/TestListenSuccess", "hello", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-reciever
	if message.PayloadString() != "hello" {
		t.Fatalf("message payload should have been 'hello' but is %v", message)
	}
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeFailure(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(context.Background())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	err = client.Subscribe(context.Background(), func(message mqtt.Message) {}, testUUID+"/#/test_publish", mqtt.ExactlyOnce) // # in the middle of a subscribe is not allowed
	if err == nil {
		t.Fatalf("subscribe should have failed: %v", err)
	}
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeSuccessAdvancedRouting(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(context.Background())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	reciever := make(chan mqtt.Message)
	err = client.Subscribe(context.Background(), func(message mqtt.Message) {
		reciever <- message
	}, testUUID+"/TestSubcribeSuccessAdvancedRouting/#", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(context.Background(), testUUID+"/TestSubcribeSuccessAdvancedRouting/abc", "hello world", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-reciever
	if message.PayloadString() != "hello world" {
		t.Fatalf("message payload should have been 'hello world' but is %v", message.PayloadString())
	}
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeNoRecieve(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	err = client.Connect(context.Background())
	defer client.DisconnectImmediately()
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	client.Listen(func(message mqtt.Message) {
		t.Fatalf("recieved a message which was not meant to happen: %v", err)
	}, testUUID+"/TestSubcribeSuccessAdvancedRouting/abc")
	err = client.Subscribe(context.Background(), nil, testUUID+"/TestSubcribeSuccessAdvancedRouting/def", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(context.Background(), testUUID+"/TestSubcribeSuccessAdvancedRouting/def", "hello world", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	<-time.After(500 * time.Millisecond)
}
