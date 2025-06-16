package consumer

import (
	"arcs/internal/configs"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Consumer struct {
	cfg configs.Config
	nc  *nats.Conn
	js  nats.JetStreamContext
}

func NewConsumerClient(cfg configs.Config) *Consumer {
	opts := []nats.Option{
		nats.Timeout(time.Duration(cfg.Consumer.ClientTimeout) * time.Second),
		nats.ReconnectWait(time.Duration(cfg.Consumer.ReconnectWait) * time.Second),
		nats.MaxReconnects(cfg.Consumer.MaxReconnects),
		//nats.ReconnectBufSize(-1), // this one disable buffer on local clients
	}

	nc, err := nats.Connect(cfg.Consumer.URL, opts...)
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to NATS: [%v]", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to jetstreams: [%v]", err)
	}

	return &Consumer{
		cfg: cfg,
		nc:  nc,
		js:  js,
	}
}

func (c *Consumer) EnsureStream() error {
	//stream alrdy exits
	_, err := c.js.StreamInfo(c.cfg.Consumer.Stream)
	if err == nil {
		return nil
	}

	_, err = c.js.AddStream(&nats.StreamConfig{
		Name:      c.cfg.Consumer.Stream,
		Subjects:  c.cfg.Consumer.Subjects,
		Retention: nats.WorkQueuePolicy,
		Storage:   nats.FileStorage,
	})
	if err != nil {
		return fmt.Errorf("[NATS] Failed to create stream [%s]: %v", c.cfg.Consumer.Stream, err)
	}

	return nil
}

func (c *Consumer) Close() {
	c.nc.Drain()
}

func (c *Consumer) Consume(topic string, handler nats.MsgHandler) error {
	if _, err := c.js.QueueSubscribe(topic, c.cfg.Consumer.Queue, handler, nats.ManualAck()); err != nil {
		return fmt.Errorf("[NATS] Failed to consume msg: %v", err)
	}

	return nil
}

func (c *Consumer) HealthCheck() error {
	if !c.nc.IsConnected() {
		return fmt.Errorf("[NATS] NATS not connected")
	}

	return nil
}
