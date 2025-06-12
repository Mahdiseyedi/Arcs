package producer

import (
	"arcs/internal/configs"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type NatsProducerClient struct {
	cfg configs.Config
	nc  *nats.Conn
	js  nats.JetStreamContext
}

func NewNatsProducerClient(cfg configs.Config) *NatsProducerClient {
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to NATS: [%v]", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to jetstreams: [%v]", err)
	}

	return &NatsProducerClient{
		cfg: cfg,
		nc:  nc,
		js:  js,
	}
}

func (c *NatsProducerClient) EnsureStream() error {
	//stream alrdy exits
	_, err := c.js.StreamInfo(c.cfg.Nats.Stream)
	if err == nil {
		return nil
	}
	
	_, err = c.js.AddStream(&nats.StreamConfig{
		Name:     c.cfg.Nats.Stream,
		Subjects: c.cfg.Nats.Subjects,
	})
	if err != nil {
		return fmt.Errorf("[NATS] Failed to create stream [%s]: %v", c.cfg.Nats.Stream, err)
	}

	return nil
}

func (c *NatsProducerClient) Publish(topic string, msg []byte) error {
	_, err := c.js.Publish(topic, msg)
	if err != nil {
		return fmt.Errorf("[NATS] Failed to publish message: %v", err)
	}
	return nil
}

func (c *NatsProducerClient) Close() {
	c.nc.Drain()
}
