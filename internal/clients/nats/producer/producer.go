package producer

import (
	"arcs/internal/configs"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Producer struct {
	cfg configs.Config
	nc  *nats.Conn
	js  nats.JetStreamContext
}

func NewProducerClient(cfg configs.Config) *Producer {
	opts := []nats.Option{
		nats.Timeout(time.Duration(cfg.Producer.ClientTimeout) * time.Second),
		nats.ReconnectWait(time.Duration(cfg.Producer.ReconnectWait) * time.Second),
		nats.MaxReconnects(cfg.Producer.MaxReconnects),
		//nats.ReconnectBufSize(-1), // this one disable buffer on local clients
	}

	nc, err := nats.Connect(cfg.Producer.URL, opts...)
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to NATS: [%v]", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[NATS] Failed to connect to jetstreams: [%v]", err)
	}

	return &Producer{
		cfg: cfg,
		nc:  nc,
		js:  js,
	}
}

func (p *Producer) EnsureStream() error {
	//stream alrdy exits
	_, err := p.js.StreamInfo(p.cfg.Producer.Stream)
	if err == nil {
		return nil
	}

	_, err = p.js.AddStream(&nats.StreamConfig{
		Name:      p.cfg.Producer.Stream,
		Subjects:  p.cfg.Producer.Subjects,
		Retention: nats.WorkQueuePolicy,
		Storage:   nats.FileStorage,
	})
	if err != nil {
		return fmt.Errorf("[NATS] Failed to create stream [%s]: %v", p.cfg.Producer.Stream, err)
	}

	return nil
}

func (p *Producer) Publish(topic string, msg []byte, idmKey string) error {
	natsMsg := &nats.Msg{
		Subject: topic,
		Data:    msg,
		Header:  nats.Header{},
	}
	natsMsg.Header.Set(nats.MsgIdHdr, idmKey)

	_, err := p.js.PublishMsg(natsMsg)
	if err != nil {
		return fmt.Errorf("[NATS] Failed to publish message: %v", err)
	}
	return nil
}

func (p *Producer) Close() {
	p.nc.Drain()
}

func (p *Producer) HealthCheck() error {
	if !p.nc.IsConnected() {
		return fmt.Errorf("[NATS] NATS not connected")
	}

	return nil
}
