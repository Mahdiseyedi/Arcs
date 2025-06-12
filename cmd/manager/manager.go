package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/repository/balance"
	"arcs/internal/repository/user"
	userSvc "arcs/internal/service/user"
	"context"
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
	var blc int64
	blc = 200

	if err := userService.CreateUser(ctx, blc); err != nil {
		log.Fatal(err)
	}
}
