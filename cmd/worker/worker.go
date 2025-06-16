package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/nats/consumer"
	"arcs/internal/configs"
	"arcs/internal/handler/worker"
	"arcs/internal/repository/sms"
	"arcs/internal/service/buffer"
	"arcs/internal/service/delivery"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	flusher := buffer.NewStatusFlusher(cfg, smsRepo)
	deliveryService := delivery.NewDeliveryService(cfg, flusher)

	//Handlers
	handler := worker.NewSMSHandler(deliveryService)

	//gracefully shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down...")
		cancel()
		flusher.Stop()
		os.Exit(0)
	}()

	if err := consumerCli.Consume(cfg.Consumer.Subjects[0], handler.Handle(ctx)); err != nil {
		log.Printf("[CONSUMER] Failed to consume message: %v", err)
	}

	select {}
}
