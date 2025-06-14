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
	validator *order.Validator
	orderSvc  *orderSvc.Svc
}

func NewOrderHandler(
	validator *order.Validator,
	svc *orderSvc.Svc) *Handler {
	return &Handler{
		validator: validator,
		orderSvc:  svc,
	}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errmsg.InvalidRequest)
		return
	}

	//TODO - add validation for order register

	if err := h.orderSvc.RegisterOrder(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}
