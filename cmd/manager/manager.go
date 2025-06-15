package main

import (
	"arcs/internal/clients/db"
	"arcs/internal/clients/nats"
	"arcs/internal/clients/redis"
	"arcs/internal/configs"
	"arcs/internal/delivery/http"
	"arcs/internal/handler/http/healthcheck"
	orderHandler "arcs/internal/handler/http/order"
	userHandler "arcs/internal/handler/http/user"
	"arcs/internal/jobs"
	balanceRepository "arcs/internal/repository/balance"
	orderRepository "arcs/internal/repository/order"
	smsRepository "arcs/internal/repository/sms"
	userRepository "arcs/internal/repository/user"
	"arcs/internal/service/health"
	orderService "arcs/internal/service/order"
	userService "arcs/internal/service/user"
	orderValidator "arcs/internal/validator/order"
	userValidator "arcs/internal/validator/user"
	"fmt"
	"time"
)

func main() {
	cfg := configs.Load("../../config.yaml")

	time.Local, _ = time.LoadLocation(cfg.Basic.TimeZone)

	//clients
	crn := jobs.NewCronJob()
	dbCli := db.NewDatabase(cfg)
	redisCli := redis.NewRedisCli(cfg)
	natsCli := nats.NewNatsClient(cfg)
	defer natsCli.Close()

	//repos
	userRepo := userRepository.NewUserRepository(dbCli)
	balanceRepo := balanceRepository.NewBalanceRepository(redisCli)
	orderRepo := orderRepository.NewOrderRepository(dbCli)
	smsRepo := smsRepository.NewSMSRepository(cfg, dbCli)

	//services
	userSvc := userService.NewUserSvc(userRepo, smsRepo, balanceRepo)
	orderSvc := orderService.NewOrderSvc(cfg, userSvc, orderRepo, smsRepo, natsCli)
	healthSvc := health.NewHealthSvc(dbCli.DB, redisCli.Client, natsCli)

	//validator
	userVal := userValidator.NewUserValidator()
	orderVal := orderValidator.NewOrderValidator()

	//handler
	healthHandle := healthcheck.NewHealthcheckHandler(healthSvc)
	userHandle := userHandler.NewUserHandler(userVal, userSvc)
	orderHandle := orderHandler.NewOrderHandler(orderVal, orderSvc)

	//jobs
	crn.C.AddFunc(fmt.Sprintf("@every %ds", cfg.Nats.RetryTimeOut), func() {
		orderSvc.RecoverUnPblishSMS()
	})

	//delivery
	server := http.NewServer(cfg, healthHandle, userHandle, orderHandle)

	crn.Start()
	server.Run()
}
