package mqtt_test

import (
	"errors"
	"testing"

	"github.com/lucacasonato/mqtt"
)

// create client with a nil server array
func TestNewClientNilServer(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{})
	if !errors.Is(err, mqtt.ErrMinimumOneServer) {
		t.Fatal("err should be ErrMinimumOneServer")
	}
	if client != nil {
		t.Fatal("client should be nil")
	}
}

// create client with a server array with no servers
func TestNewClientNoServer(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{},
	})
	if !errors.Is(err, mqtt.ErrMinimumOneServer) {
		t.Fatal("err should be ErrMinimumOneServer")
	}
	if client != nil {
		t.Fatal("client should be nil")
	}
}

// create client with a server array with no servers
func TestNewClientBasicServer(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatal("err should be nil")
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
}

// check that a client gets created and a client id is generated when it is not set
func TestNewClientNoClientID(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
	})
	if err != nil {
		t.Fatalf("err should be nil but is %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.Options.ClientID == "" {
		t.Fatal("client.Options.ClientID should not be empty")
	}
}

// check that a client gets created and a client id is not changed when it is set
func TestNewClientHasClientID(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1883",
		},
		ClientID: "client-id",
	})
	if err != nil {
		t.Fatal("err should be nil")
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.Options.ClientID != "client-id" {
		t.Fatal("client.Options.ClientID should be 'client-id'")
	}

}
