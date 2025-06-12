package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/repository/user"
	"log"
)

func main() {
	cfg := configs.Load("../../config.yaml")

	//clients
	dbClient := db.NewDatabase(cfg)
	redisClient := redis.NewRedisCli(cfg)

	//repos
	userRepo := user.NewUserRepository(dbClient)
	log.Println(redisClient)
	log.Println(userRepo)
}
