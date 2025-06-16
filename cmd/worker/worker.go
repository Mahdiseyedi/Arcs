package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/nats/consumer"
	"arcs/internal/configs"
	"arcs/internal/handler/worker"
	"arcs/internal/repository/sms"
	"arcs/internal/service/delivery"
	"context"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configs.Load("../../worker-config.yaml")
	time.Local, _ = time.LoadLocation(cfg.Basic.TimeZone)

	//clients
	dbCli := db.NewDatabase(cfg)
	consumerCli := consumer.NewConsumerClient(cfg)
	defer consumerCli.Close()

	//repositories
	smsRepo := sms.NewSMSRepository(cfg, dbCli)

	//service
	deliveryService := delivery.NewDeliveryService(cfg, smsRepo)

	//Handlers
	handler := worker.NewSMSHandler(deliveryService)

	for {
		if err := consumerCli.Consume(cfg.Consumer.Subjects[0], handler.Handle(ctx)); err != nil {
			log.Printf("[CONSUMER] Failed to consume message: %v", err)
			time.Sleep(1 * time.Second)
		}
	}
}
