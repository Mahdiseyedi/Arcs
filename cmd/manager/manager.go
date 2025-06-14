package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/nats"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/delivery/http"
	"arcs/internal/dto"
	"arcs/internal/handler/http/healthcheck"
	orderHandler "arcs/internal/handler/http/order"
	userHandler "arcs/internal/handler/http/user"
	balanceRepository "arcs/internal/repository/balance"
	orderRepository "arcs/internal/repository/order"
	userRepository "arcs/internal/repository/user"
	orderService "arcs/internal/service/order"
	userService "arcs/internal/service/user"
	orderValidator "arcs/internal/validator/order"
	userValidator "arcs/internal/validator/user"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	cfg := configs.Load("../../config.yaml")

	//clients
	dbCli := db.NewDatabase(cfg)
	redisCli := redis.NewRedisCli(cfg)
	natsCli := nats.NewNatsClient(cfg)

	//repos
	userRepo := userRepository.NewUserRepository(dbCli)
	balanceRepo := balanceRepository.NewBalanceRepository(redisCli)
	orderRepo := orderRepository.NewOrderRepository(dbCli)

	//services
	userSvc := userService.NewUserSvc(userRepo, balanceRepo)
	orderSvc := orderService.NewOrderSvc(cfg, userSvc, orderRepo, natsCli)

	//validator
	userVal := userValidator.NewUserValidator()
	orderVal := orderValidator.NewOrderValidator()

	//handler
	healthHandle := healthcheck.NewHealthcheckHandler()
	userHandle := userHandler.NewUserHandler(userVal, userSvc)
	orderHandle := orderHandler.NewOrderHandler(orderVal, orderSvc)

	server := http.NewServer(cfg, healthHandle, userHandle, orderHandle)

	server.Run()

	time.Sleep(time.Minute * 100)

	//TODO - remove me
	ctx := context.Background()
	dummyID := "01c0cad0-ebea-43ee-a92f-3230d00b4f0e"

	defer natsCli.Close()

	userSvc.CreateUser(ctx, dto.CreateUserRequest{Balance: 200})

	d, err := userSvc.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("before order Balance: ", d)

	if err := orderSvc.RegisterOrder(ctx, dto.OrderRequest{
		UserID:  dummyID,
		Content: "hi guys",
		Destinations: []string{
			"0912",
			"0913",
			"0914",
			"0915",
			"0916",
			"0917",
			"0918",
			"0919",
		},
	}); err != nil {
		log.Fatal(err)
	}

	t, err := userSvc.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("after order Balance: ", t)
}
