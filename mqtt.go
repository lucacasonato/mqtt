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
	router  *router
}

// ClientOptions is the list of options used to create a client
type ClientOptions struct {
	Servers  []string // The list of broker hostnames to connect to
	ClientID string   // If left empty a uuid will automatically be generated
	Username string   // If not set then authentication will not be used
	Password string   // Will only be used if the username is set

	AutoReconnect bool // If the client should automatically try to reconnect when the connection is lost
}

// QOS describes the quality of service of an mqtt publish
type QOS byte

const (
	AtMostOnce  QOS = iota // Deliver at most once to every subscriber - this means message delivery is not guaranteed
	AtLeastOnce            // Deliver a message at least once to every subscriber
	ExactlyOnce            // Deliver a message exactly once to every subscriber
)

var (
	ErrMinimumOneServer = errors.New("mqtt: at least one server needs to be specified")
)

func handle(callback MessageHandler) paho.MessageHandler {
	return func(client paho.Client, message paho.Message) {
		if callback != nil {
			callback(Message{message: message})
		}
	}
}

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
	router := newRouter()
	pahoClient.AddRoute("#", handle(func(message Message) {
		routes := router.match(&message)
		for _, route := range routes {
			route.handler(message)
		}
	}))

	return &Client{client: pahoClient, Options: options, router: router}, nil
}

// Connect tries to establish a conenction with the mqtt servers
func (c *Client) Connect(ctx context.Context) error {
	// try to connect to the client
	token := c.client.Connect()
	return tokenWithContext(ctx, token)
}

// Disconnect will immediately close the conenction with the mqtt servers
func (c *Client) DisconnectImmediately() {
	c.client.Disconnect(0)
}

func tokenWithContext(ctx context.Context, token paho.Token) error {
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
