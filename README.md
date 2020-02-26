# mqtt

[![GoDoc](https://godoc.org/github.com/lucacasonato/mqtt?status.svg)](http://godoc.org/github.com/lucacasonato/mqtt)
[![CI](https://github.com/lucacasonato/mqtt/workflows/ci/badge.svg)](https://github.com/lucacasonato/mqtt/actions?workflow=ci)
[![Code Coverage](https://img.shields.io/codecov/c/gh/lucacasonato/mqtt)](https://codecov.io/gh/lucacasonato/mqtt)
[![Go Report](https://goreportcard.com/badge/github.com/lucacasonato/mqtt)](https://goreportcard.com/report/github.com/lucacasonato/mqtt)

An mqtt client for Go that improves usability over the [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) library it wraps. Made for 🧑.

## installation

```bash
go get github.com/lucacasonato/mqtt
```

```go
import "github.com/lucacasonato/mqtt"
// or
import (
    "github.com/lucacasonato/mqtt"
)
```

## usage

### creating a client & connecting

```go
client, err := mqtt.NewClient(mqtt.ClientOptions{
    // required
    Servers: []string{
        "tcp://test.mosquitto.org:1883",
    },

    // optional
    ClientID: "my-mqtt-client",
    Username: "admin",
    Password: "***",
    AutoReconnect: true,
})
if err != nil {
    panic(err)
}

err = client.Connect(context.WithTimeout(2 * time.Second))
if err != nil {
    panic(err)
}
```

You can use any of these schemes for the broker `tcp` (unesecured), `ssl` (secured), `ws` (unsecured), `wss` (secured).

### disconnecting from a client

```go
client.Disconnect()
```

### publishing a message

#### bytes

```go
err := client.Publish(context.WithTimeout(1 * time.Second), "api/v0/main/client1", []byte(0, 1 ,2, 3), mqtt.AtLeastOnce)
if err != nil {
    panic(err)
}
```

#### string

```go
err := client.PublishString(context.WithTimeout(1 * time.Second), "api/v0/main/client1", "hello world", mqtt.AtLeastOnce)
if err != nil {
    panic(err)
}
```

#### json

```go
err := client.PublishJSON(context.WithTimeout(1 * time.Second), "api/v0/main/client1", []string("hello", "world"), mqtt.AtLeastOnce)
if err != nil {
    panic(err)
}
```

### subscribing

```go
err := client.Subscribe(context.WithTimeout(1 * time.Second), "api/v0/main/client1", mqtt.AtLeastOnce)
if err != nil {
    panic(err)
}
```

```go
err := client.SubscribeMultiple(context.WithTimeout(1 * time.Second), map[string]mqtt.QOS{
    "api/v0/main/client1": mqtt.AtLeastOnce,
})
if err != nil {
    panic(err)
}
```

### handling

```go
route := client.Handle("api/v0/main/client1", func(message mqtt.Message) {
    v := interface{}{}
    err := message.PayloadJSON(&v)
    if err != nil {
        panic(err)
    }
    fmt.Printf("recieved a message with content %v\n", v)
})
// once you are done with the route you can stop handling it
route.Stop()
```

### listening

```go
messages, route := client.Listen("api/v0/main/client1")
for {
    message := <-messages
    fmt.Printf("recieved a message with content %v\n", message.PayloadString())
}
// once you are done with the route you can stop handling it
route.Stop()
```
