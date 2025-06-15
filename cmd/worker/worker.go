package main

import (
	natsCli "arcs/internal/clients/nats"
	"arcs/internal/configs"
	"arcs/internal/models"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	cfg := configs.Load("../../worker-config.yaml")
	time.Local, _ = time.LoadLocation(cfg.Basic.TimeZone)

	natsClient := natsCli.NewNatsClient(cfg)
	_ = natsClient.EnsureStream()

	for {
		if err := natsClient.Consume(cfg.Nats.Subjects[0], createMessageHandler()); err != nil {
			log.Printf("[CONSUMER] Failed to consume message: %v", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func createMessageHandler() nats.MsgHandler {
	return func(msg *nats.Msg) {
		var sms models.SMS
		// TODO - replace me with protobuf
		if err := json.Unmarshal(msg.Data, &sms); err != nil {
			log.Printf("Failed to unmarshal msg: %v", err)
			return
		}

		time.Sleep(1 * time.Second) // simulate job

		msg.Ack() // manual ack

		log.Printf("Processed: [%s] - [%s] - [%s]\n", sms.ID, sms.Destination, sms.Order.Content)
	}
}
