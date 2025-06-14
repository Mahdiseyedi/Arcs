package order

import (
	"arcs/internal/dto"
	orderSvc "arcs/internal/service/order"
	"arcs/internal/utils/errmsg"
	"arcs/internal/validator/order"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	order.Validator
	orderSvc *orderSvc.Svc
}

func NewOrderHandler(orderSvc *orderSvc.Svc) *Handler {
	return &Handler{
		orderSvc: orderSvc,
	}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errmsg.InvalidRequest)
		return
	}

	if err := h.orderSvc.RegisterOrder(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}
