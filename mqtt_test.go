package mqtt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lucacasonato/mqtt"
)

// TestNewClientNilServer checks if creating a client with a nil server array works
func TestNewClientNilServer(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{})
	if !errors.Is(err, mqtt.ErrMinimumOneServer) {
		t.Fatal("err should be ErrMinimumOneServer")
	}
	if client != nil {
		t.Fatal("client should be nil")
	}
}

// TestNewClientNoServer checks if creating a client with a server array with no servers works
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

// TestNewClientBasicServer checks if creating a client with a server array with one server works
func TestNewClientBasicServer(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
}

// TestNewClientNoClientID checks that a client gets created and a client id is generated when it is not set
func TestNewClientNoClientID(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.Options.ClientID == "" {
		t.Fatal("client.Options.ClientID should not be empty")
	}
}

// TestNewClientHasClientID checks that a client gets created and a client id is not changed when it is already set
func TestNewClientHasClientID(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
		ClientID: "client-id",
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	if client.Options.ClientID != "client-id" {
		t.Fatal("client.Options.ClientID should be 'client-id'")
	}
}

// TestNewClientWithAuthentication has username and password to check if those get set
func TestNewClientWithAuthentication(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
		Username: "user",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
}

// TestConnectSuccess just checks that connecting to a broker works
func TestConnectSuccess(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	err = client.Connect(ctx())
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
}

// TestConnectContextTimeout checks if connect errors if a context with a timeout times out
func TestConnectContextTimeout(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	ctx, cancel := context.WithTimeout(ctx(), 1*time.Nanosecond)
	defer cancel()
	err = client.Connect(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("connect should have failed with error context.DeadlineExceeded")
	}
}

// TestConnectContextCancel checks if connect errors if a context with a cancel gets canceled
func TestConnectContextCancel(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	ctx, cancel := context.WithCancel(ctx())
	go func() {
		time.Sleep(1 * time.Microsecond)
		cancel()
	}()
	defer cancel()
	err = client.Connect(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatal("connect should have failed with error context.Canceled")
	}
}

// TestConnectFailed that a invalid client does not connect and errors
func TestConnectFailed(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			"tcp://test.mosquitto.org:1884", // incorrect port
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	err = client.Connect(ctx())
	if err == nil {
		t.Fatal("connect should have failed")
	}
}

// TestDisconnectImmediately immediately disconnects the mqtt broker
func TestDisconnectImmediately(t *testing.T) {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{
			broker,
		},
	})
	if err != nil {
		t.Fatalf("creating client failed: %v", err)
	}
	if client == nil {
		t.Fatal("client should not be nil")
	}
	err = client.Connect(ctx())
	if err != nil {
		t.Fatalf("connect should not have failed: %v", err)
	}
	client.DisconnectImmediately()
}
