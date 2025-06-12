package main

import (
	natsCli "arcs/internal/clients/nats"
	"arcs/internal/configs"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	cfg := configs.Load("../../worker-config.yaml")

	natsClient := natsCli.NewNatsClient(cfg)

	for {
		if err := natsClient.Consume(cfg.Nats.Subjects[0], createMessageHandler("1")); err != nil {
			log.Printf("[CONSUMER] Failed to consume message: %v", err)
		}
	}
}

func createMessageHandler(workerID string) nats.MsgHandler {
	return func(msg *nats.Msg) {
		log.Printf("[%s] Processing: %s\n\n", workerID, string(msg.Data))

		time.Sleep(1 * time.Second) // simulate job

		msg.Ack() // manual ack

		log.Printf("[%s] Done: %s\n\n", workerID, string(msg.Data))
	}
}
