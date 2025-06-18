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
	r := gin.New()

	//health
	r.GET("/health", s.heathHandler.Check)
	r.Use(gin.Logger(), gin.Recovery())

	//api version 1
	api := r.Group("/api/v1")
	{
		u := api.Group("/user")
		{
			u.POST("/", s.userHandler.CreateUser)
			u.POST("/charge", s.userHandler.ChargeUser)
			u.GET("/balance/:id", s.userHandler.GetUserBalance)
			u.GET("/:id/sms", s.userHandler.GetFilteredUserSMS)
		}

		o := api.Group("/order")
		{
			o.POST("/", s.orderHandler.CreateOrder)
		}
	}

	port := fmt.Sprintf(":%s", s.cfg.Basic.Port)
	if err := r.Run(port); err != nil {
		log.Fatalf("[ROUTER] Http server start error: %v", err)
	}
}
