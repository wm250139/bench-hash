package main

import (
	"bench-hash/hasher"
	"context"
	"flag"
	"fmt"
	mq "github.com/eclipse/paho.golang/paho"
	"log"
	"net"
)

var (
	netType = flag.String("type", "tcp", "Network type for net.Dial")
	addr    = flag.String("addr", "localhost:1883", "The connection address for the MQTT broker")
)

func init() {
	flag.Parse()
}

func main() {
	// Connect to the broker
	conn, err := net.Dial(*netType, *addr)
	if err != nil {
		log.Fatal("unable to reach broker", err)
	}

	// Create client instance
	client := mq.NewClient()
	client.Conn = conn
	client.Router.RegisterHandler("hasher/request", func(p *mq.Publish) {
		hash := hasher.String(string(p.Payload))
		fmt.Printf("Got request for hashing '%s', calculated '%s'\nResponse topic: %s\n", string(p.Payload), hash, p.Properties.ResponseTopic)
		// Ignoring response as we don't need it for this test
		_, e := client.Publish(context.Background(), &mq.Publish{
			QoS:        0,
			Topic:      p.Properties.ResponseTopic,
			Payload:    []byte(hash),
			Properties: &mq.PublishProperties{CorrelationData: p.Properties.CorrelationData},
		})
		if e != nil {
			log.Fatal("unable to send response", e)
		}
	})

	// Connect options
	connectOpts := &mq.Connect{ClientID: "go-srv", KeepAlive: 30, CleanStart: true}

	// Ignoring ack packet
	_, err = client.Connect(context.Background(), connectOpts)
	if err != nil {
		log.Fatalf("unable to connect client: %s", err)
	}

	// Ignoring ack packet
	_, err = client.Subscribe(context.Background(), &mq.Subscribe{
		Subscriptions: map[string]mq.SubscribeOptions{
			"hasher/request": {QoS: 0},
		},
	})
	if err != nil {
		log.Fatalf("Unable to subscribe: %s", err)
	}

	fmt.Println("Server started. Press enter to stop.")
	var input string
	fmt.Scanln(&input)
	fmt.Println("Server stopping.")
	_ = client.Disconnect(&mq.Disconnect{})
}
