package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/nats"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/dto"
	"arcs/internal/repository/balance"
	"arcs/internal/repository/order"
	"arcs/internal/repository/user"
	orderSvc "arcs/internal/service/order"
	userSvc "arcs/internal/service/user"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := configs.Load("../../config.yaml")

	//clients
	dbClient := db.NewDatabase(cfg)
	redisClient := redis.NewRedisCli(cfg)
	natsClient := nats.NewNatsClient(cfg)

	//repos
	userRepo := user.NewUserRepository(dbClient)
	balanceRepo := balance.NewBalanceRepository(redisClient)
	orderRepo := order.NewOrderRepository(dbClient)

	//services
	userService := userSvc.NewUserSvc(userRepo, balanceRepo)
	orderService := orderSvc.NewOrderSvc(cfg, userService, orderRepo, natsClient)

	//TODO - remove me
	ctx := context.Background()
	dummyID := "01c0cad0-ebea-43ee-a92f-3230d00b4f0e"

	defer natsClient.Close()

	userService.CreateUser(ctx, dto.CreateUserRequest{Balance: 200})

	d, err := userService.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("before order Balance: ", d)

	if err := orderService.RegisterOrder(ctx, dto.OrderRequest{
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

	t, err := userService.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("after order Balance: ", t)
}
