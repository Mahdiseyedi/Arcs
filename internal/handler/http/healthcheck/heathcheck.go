package healthcheck

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
}

func NewHealthcheckHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Check(c *gin.Context) {
	//TODO - add more logic for health check here
	c.Status(http.StatusOK)
}
