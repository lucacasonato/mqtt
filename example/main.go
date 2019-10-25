package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lucacasonato/mqtt"
)

func ctx() context.Context {
	cntx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	return cntx
}

type Color struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

func main() {
	client, err := mqtt.NewClient(mqtt.ClientOptions{
		Servers: []string{"tcp://localhost:1883"},
	})
	if err != nil {
		log.Fatalf("failed to create mqtt client: %v\n", err)
	}

	err = client.Connect(ctx())
	if err != nil {
		log.Fatalf("failed to connect to mqtt server: %v\n", err)
	}

	err = client.Subscribe(ctx(), "my-home-automation/lamps/#", mqtt.AtMostOnce)
	if err != nil {
		log.Fatalf("failed to subscribe to config service: %v\n", err)
	}

	client.Handle("my-home-automation/lamps/+/color", func(m mqtt.Message) {
		lampID := m.TopicVars()[0]
		var color Color
		err := m.PayloadJSON(&color)
		if err != nil {
			log.Printf("failed to parse color: %v\n", err)
			return
		}
		log.Printf("lamp %v now has the color r: %v g: %v b: %v\n", lampID, color.Red, color.Blue, color.Green)
	})

	for {
		lampID := uuid.New().String()
		err := client.PublishJSON(ctx(), "my-home-automation/lamps/"+lampID+"/color", Color{
			Red:   uint8(rand.Intn(255)),
			Green: uint8(rand.Intn(255)),
			Blue:  uint8(rand.Intn(255)),
		}, mqtt.AtLeastOnce)
		if err != nil {
			log.Printf("failed to publish: %v\n", err)
			continue
		}

		<-time.After(1 * time.Second)
	}
}
