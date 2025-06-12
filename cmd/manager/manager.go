package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/repository/balance"
	"arcs/internal/repository/user"
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

	//repos
	userRepo := user.NewUserRepository(dbClient)
	balanceRepo := balance.NewBalanceRepository(redisClient)

	//services
	userService := userSvc.NewUserSvc(userRepo, balanceRepo)
	//TODO - remove me
	ctx := context.Background()
	dummyID := "8de01a14-9532-4ee1-af82-badb92dfe7da"

	v, err := userService.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("before charge: [%v]\n", v)

	if err := userService.ChargeUser(ctx, dummyID, 33); err != nil {
		log.Fatal(err)
	}

	d, err := userService.Balance(ctx, dummyID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("after charge: [%v]\n", d)
}
