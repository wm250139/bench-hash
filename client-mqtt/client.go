package client_mqtt

import (
	"context"
	"errors"
	"fmt"
	mq "github.com/eclipse/paho.golang/paho"
	"net"
	"sync"
	"time"
)

type MQTTClient struct {
	m         sync.Mutex
	client    *mq.Client
	listeners map[string]chan string
}

func NewMQTTClient(netType string, addr string) (*MQTTClient, error) {
	// Connect to the broker
	conn, err := net.Dial(netType, addr)
	if err != nil {
		return nil, err
	}

	// Create client instance
	client := mq.NewClient()
	client.Conn = conn

	// Connect options
	connectOpts := &mq.Connect{ClientID: "go-client", KeepAlive: 30, CleanStart: true}

	// Ignoring ack packet
	_, err = client.Connect(context.Background(), connectOpts)
	if err != nil {
		return nil, err
	}

	//client.OnDisconnect = func(d packets.Disconnect) {
	//	_, _ = client.Connect(context.Background(), connectOpts)
	//}

	listeners := make(map[string]chan string)

	client.Router.RegisterHandler(fmt.Sprintf("%s/responses", client.ClientID), func(p *mq.Publish) {
		corID := string(p.Properties.CorrelationData)
		s := string(p.Payload)
		listeners[corID] <- s
	})

	// Subscribe for response
	_, err = client.Subscribe(context.Background(), &mq.Subscribe{
		Subscriptions: map[string]mq.SubscribeOptions{
			fmt.Sprintf("%s/responses", client.ClientID): {QoS: 0},
		},
	})
	if err != nil {
		return nil, errors.New("unable to subscribe for responses")
	}

	return &MQTTClient{client: client, listeners: listeners}, nil
}

func (c *MQTTClient) Hash(input string) (string, error) {
	corID := fmt.Sprintf("%s", time.Now())
	respChan := make(chan string, 1)
	c.addCorrelationListener(corID, respChan)

	// Publish the request
	_, err := c.client.Publish(context.Background(), &mq.Publish{
		QoS:     0,
		Topic:   "hasher/request",
		Payload: []byte(input),
		Properties: &mq.PublishProperties{
			ResponseTopic:   fmt.Sprintf("%s/responses", c.client.ClientID),
			CorrelationData: []byte(corID),
		},
	})
	if err != nil {
		return "", errors.New("unable to send request")
	}

	resp := <-respChan

	return resp, nil
}

func (c *MQTTClient) addCorrelationListener(id string, channel chan string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.listeners[id] = channel
}

func (c *MQTTClient) Close() {
	c.client.Disconnect(&mq.Disconnect{})
}
