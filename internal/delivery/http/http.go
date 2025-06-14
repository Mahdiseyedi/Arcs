package http

import (
	"arcs/internal/configs"
	"arcs/internal/handler/http/healthcheck"
	"arcs/internal/handler/http/order"
	"arcs/internal/handler/http/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	cfg          configs.Config
	heathHandler *healthcheck.Handler
	userHandler  *user.Handler
	orderHandler *order.Handler
}

func NewServer(
	cfg configs.Config,
	healthHandler *healthcheck.Handler,
	userHandler *user.Handler,
	orderHandler *order.Handler,
) *Server {
	return &Server{
		cfg:          cfg,
		heathHandler: healthHandler,
		userHandler:  userHandler,
		orderHandler: orderHandler,
	}
}

func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//health
	r.GET("/health", s.heathHandler.Check)

	//user
	r.POST("api/v1/user", s.userHandler.CreateUser)
	r.POST("api/v1/user/charge", s.userHandler.ChargeUser)
	r.GET("/api/v1/user/balance/:id", s.userHandler.GetUserBalance)

	//order
	r.POST("/api/v1/order", s.orderHandler.CreateOrder)

	//add more routes here

	port := fmt.Sprintf(":%s", s.cfg.Basic.Port)
	if err := r.Run(port); err != nil {
		log.Fatalf("[ROUTER] Http server start error: %v", err)
	}
}
