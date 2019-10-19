package mqtt

import (
	"context"
	"errors"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// Client for talking using mqtt
type Client struct {
	Options ClientOptions // The options that were used to create this client
	client  paho.Client
}

// ClientOptions is the list of options used to create a client
type ClientOptions struct {
	Servers  []string // The list of broker hostnames to connect to
	ClientID string   // If left empty a uuid will automatically be generated
	Username string   // If not set then authentication will not be used
	Password string   // Will only be used if the username is set

	AutoReconnect bool // If the client should automatically try to reconnect when the connection is lost
}

var (
	ErrMinimumOneServer = errors.New("mqtt: at least one server needs to be specified")
)

// NewClient creates a new client with the specified options
func NewClient(options ClientOptions) (*Client, error) {
	pahoOptions := paho.NewClientOptions()

	// brokers
	if options.Servers != nil && len(options.Servers) > 0 {
		for _, server := range options.Servers {
			pahoOptions.AddBroker(server)
		}
	} else {
		return nil, ErrMinimumOneServer
	}

	// client id
	if options.ClientID == "" {
		options.ClientID = uuid.New().String()
	}
	pahoOptions.SetClientID(options.ClientID)

	// auth
	if options.Username != "" {
		pahoOptions.SetUsername(options.Username)
		pahoOptions.SetPassword(options.Password)
	}

	// auto reconnect
	pahoOptions.SetAutoReconnect(options.AutoReconnect)

	pahoClient := paho.NewClient(pahoOptions)
	return &Client{client: pahoClient, Options: options}, nil
}

// Connect tries to establish a conenction with the mqtt servers
func (c *Client) Connect(ctx context.Context) error {
	// try to connect to the client
	token := c.client.Connect()
	completer := make(chan error)

	// TODO: This go routine will not be removed up if the ctx is cancelled or a the ctx timeout passes
	go func() {
		token.Wait()
		completer <- token.Error()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-completer:
			return err
		}
	}
}
