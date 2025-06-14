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
	cfg configs.Config
	hc  *healthcheck.Handler

	u *user.Handler
	o *order.Handler
}

func NewServer(
	cfg configs.Config,
	hc *healthcheck.Handler,
	u *user.Handler,
	o *order.Handler,
) *Server {
	return &Server{
		cfg: cfg,
		hc:  hc,
		u:   u,
		o:   o,
	}
}

func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/health", s.hc.Check)

	//user
	r.GET("/api/v1/user/balance/:id", s.u.GetUserBalance)

	//order
	r.POST("/api/v1/order", s.o.CreateOrder)

	//add more routes here

	port := fmt.Sprintf(":%s", s.cfg.Basic.Port)
	if err := r.Run(port); err != nil {
		log.Fatalf("[ROUTER] Http server start error: %v", err)
	}
}
