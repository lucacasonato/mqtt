package mqtt_test

import (
	"context"
	"testing"
	"time"

	"github.com/lucacasonato/mqtt"
)

func ctx() context.Context {
	c, _ := context.WithTimeout(context.Background(), 1*time.Second)
	return c
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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

	receiver := make(chan mqtt.Message)
	err = client.Subscribe(ctx(), testUUID+"/TestSubcribeSuccess/#", mqtt.ExactlyOnce)
	client.Handle(testUUID+"/TestSubcribeSuccess/#", func(message mqtt.Message) {
		receiver <- message
	})
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestSubcribeSuccess/abc", "[1, 2]", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-receiver
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
	if message.Topic() != testUUID+"/TestSubcribeSuccess/abc" {
		t.Fatalf("message topic should be %v but is %v", testUUID+"/TestSubcribeSuccess/abc", message.Topic())
	}
	if message.QOS() != mqtt.ExactlyOnce {
		t.Fatalf("message qos should be mqtt.ExactlyOnce but is %v", message.QOS())
	}
	if message.IsDuplicate() != false {
		t.Fatalf("message IsDuplicate should be false but is %v", message.IsDuplicate())
	}
	vars := message.TopicVars()
	if len(vars) != 1 && vars[0] != "abc" {
		t.Fatalf("message TopicVars should be ['abc'] but is %v", vars)
	}

	message.Acknowledge()
}

// TestSubcribeMultipleSuccess checks that a message gets recieved correctly
func TestSubcribeMultipleSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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

	receiver := make(chan mqtt.Message)
	err = client.SubscribeMultiple(ctx(), map[string]mqtt.QOS{testUUID + "/TestSubcribeMultipleSuccess/#": mqtt.ExactlyOnce})
	client.Handle(testUUID+"/TestSubcribeMultipleSuccess/#", func(message mqtt.Message) {
		receiver <- message
	})
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestSubcribeMultipleSuccess/abc", "[1, 2]", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-receiver
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
	if message.Topic() != testUUID+"/TestSubcribeMultipleSuccess/abc" {
		t.Fatalf("message topic should be %v but is %v", testUUID+"/TestSubcribeMultipleSuccess/abc", message.Topic())
	}
	if message.QOS() != mqtt.ExactlyOnce {
		t.Fatalf("message qos should be mqtt.ExactlyOnce but is %v", message.QOS())
	}
	if message.IsDuplicate() != false {
		t.Fatalf("message IsDuplicate should be false but is %v", message.IsDuplicate())
	}
	vars := message.TopicVars()
	if len(vars) != 1 && vars[0] != "abc" {
		t.Fatalf("message TopicVars should be ['abc'] but is %v", vars)
	}

	message.Acknowledge()
}

// TestListenSuccess checks that a listener recieves a message correctly
func TestListenSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	receiver, _ := client.Listen(testUUID + "/TestListenSuccess")
	err = client.Subscribe(ctx(), testUUID+"/TestListenSuccess", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestListenSuccess", "hello", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-receiver
	if message.PayloadString() != "hello" {
		t.Fatalf("message payload should have been 'hello' but is %v", message)
	}
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeFailure(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	err = client.Subscribe(ctx(), testUUID+"/#/test_publish", mqtt.ExactlyOnce) // # in the middle of a subscribe is not allowed
	if err == nil {
		t.Fatalf("subscribe should have failed: %v", err)
	}
}

// TestSubcribeSuccess checks that a message gets recieved correctly
func TestSubcribeSuccessAdvancedRouting(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	receiver := make(chan mqtt.Message)
	client.Handle(testUUID+"/TestSubcribeSuccessAdvancedRouting/#", func(message mqtt.Message) {
		receiver <- message
	})
	err = client.Subscribe(ctx(), testUUID+"/TestSubcribeSuccessAdvancedRouting/#", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestSubcribeSuccessAdvancedRouting/abc", "hello world", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	message := <-receiver
	if message.PayloadString() != "hello world" {
		t.Fatalf("message payload should have been 'hello world' but is %v", message.PayloadString())
	}
}

// TestSubcribeNoRecieve checks that a message does not get recieved when it is not listening
func TestSubcribeNoRecieve(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	client.Handle(testUUID+"/TestSubcribeNoRecieve/abc", func(message mqtt.Message) {
		t.Fatalf("recieved a message which was not meant to happen: %v", err)
	})
	err = client.Subscribe(ctx(), testUUID+"/TestSubcribeNoRecieve/def", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestSubcribeNoRecieve/def", "hello world", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	<-time.After(500 * time.Millisecond)
}

// TestUnsubcribe checks that a message does not get recieved after you unsubscribe
func TestUnsubcribe(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	client.Handle(testUUID+"/TestUnsubcribe", func(message mqtt.Message) {
		t.Fatalf("recieved a message which was not meant to happen: %v", err)
	})
	err = client.Subscribe(ctx(), testUUID+"/TestUnsubcribe", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.Unsubscribe(ctx(), testUUID+"/TestUnsubcribe")
	if err != nil {
		t.Fatalf("unsubscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestUnsubcribe", "hello world", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	<-time.After(500 * time.Millisecond)
}

// TestRemoveRoute checks that a route can be unsubscribed from
func TestRemoveRoute(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
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
	receiver, route := client.Listen(testUUID + "/TestRemoveRoute")
	err = client.Subscribe(ctx(), testUUID+"/TestRemoveRoute", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("subscribe should not have failed: %v", err)
	}
	err = client.PublishString(ctx(), testUUID+"/TestRemoveRoute", "hello", mqtt.ExactlyOnce)
	if err != nil {
		t.Fatalf("publish should not have failed: %v", err)
	}
	<-receiver
	route.Stop()
	select {
	case <-receiver:
		t.Fatalf("recieved a message which was not meant to happen: %v", err)
	case <-time.After(500 * time.Millisecond):
	}
}

// TestEmptyRoute checks that an empty route does nothing
func TestEmptyRoute(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client should not have failed: %v", err)
	}
	client.Handle(testUUID+"/TestEmptyRoute/abc", nil)
}
