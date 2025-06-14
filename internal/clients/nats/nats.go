package nats

import (
	"arcs/internal/configs"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type Client struct {
	cfg configs.Config
	nc  *nats.Conn
	js  nats.JetStreamContext
}

func NewNatsClient(cfg configs.Config) *Client {
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to NATS: [%v]", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to jetstreams: [%v]", err)
	}

	return &Client{
		cfg: cfg,
		nc:  nc,
		js:  js,
	}
}

func (c *Client) EnsureStream() error {
	//stream alrdy exits
	_, err := c.js.StreamInfo(c.cfg.Nats.Stream)
	if err == nil {
		return nil
	}

	_, err = c.js.AddStream(&nats.StreamConfig{
		Name:      c.cfg.Nats.Stream,
		Subjects:  c.cfg.Nats.Subjects,
		Retention: nats.WorkQueuePolicy,
		Storage:   nats.FileStorage,
	})
	if err != nil {
		return fmt.Errorf("[NATS] Failed to create stream [%s]: %v", c.cfg.Nats.Stream, err)
	}

	return nil
}

func (c *Client) Publish(topic string, msg []byte) error {
	_, err := c.js.Publish(topic, msg)
	if err != nil {
		return fmt.Errorf("[NATS] Failed to publish message: %v", err)
	}
	return nil
}

func (c *Client) Close() {
	c.nc.Drain()
}

func (c *Client) Consume(topic string, handler nats.MsgHandler) error {
	if _, err := c.js.QueueSubscribe(topic, c.cfg.Nats.Queue, handler, nats.ManualAck()); err != nil {
		return fmt.Errorf("[NATS] Failed to consume msg: %v", err)
	}

	return nil
}

func (c *Client) HealthCheck() error {
	if !c.nc.IsConnected() {
		return fmt.Errorf("[NATS] NATS not connected")
	}

	return nil
}
