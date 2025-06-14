package healthcheck

import (
	"arcs/internal/service/health"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	svc *health.Svc
}

func NewHealthcheckHandler(
	svc *health.Svc,
) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Check(c *gin.Context) {
	if err := h.svc.DBHealthCheck(c.Request.Context()); err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	if err := h.svc.RedisHealthCheck(c.Request.Context()); err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	if err := h.svc.NatsHealthCheck(c.Request.Context()); err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	//TODO - add more logic for health check here

	c.Status(http.StatusOK)
}
