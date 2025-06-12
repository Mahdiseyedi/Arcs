package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/repository/balance"
	"arcs/internal/repository/user"
	"context"
	"github.com/google/uuid"
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

	//TODO - remove me
	ctx := context.Background()
	var blc int64
	blc = 200

	id, err := userRepo.CreateUser(ctx, uuid.NewString(), blc)
	if err != nil {
		log.Fatal(err)
	}

	if err := balanceRepo.Set(ctx, id, blc); err != nil {
		log.Fatal(err)
	}

	bc, err := balanceRepo.Get(ctx, id)
	log.Printf("user: [%v], balance: [%v]", id, bc)
}
