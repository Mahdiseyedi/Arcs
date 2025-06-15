package user

import (
	"arcs/internal/dto"
	"arcs/internal/service/user"
	"arcs/internal/utils"
	"arcs/internal/utils/errmsg"
	validator "arcs/internal/validator/user"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	validator *validator.Validator
	userSvc   *user.Svc
}

func NewUserHandler(
	validator *validator.Validator,
	userSvc *user.Svc,
) *Handler {
	return &Handler{
		validator: validator,
		userSvc:   userSvc,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg.InvalidRequest.Error()})
		return
	}

	if err := h.validator.CreateUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userSvc.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": errmsg.FailedCreateUser.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ChargeUser(c *gin.Context) {
	var req dto.ChargeUserBalance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg.InvalidRequest.Error()})
		return
	}

	if err := h.validator.ChargeUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userSvc.ChargeUser(c.Request.Context(), req); err != nil {
		if errors.Is(err, errmsg.UserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": errmsg.UserNotFound.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to charge user balance"})
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *Handler) GetUserBalance(c *gin.Context) {
	userID := c.Param("id")
	if !utils.IsValidUUID(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg.InvalidID.Error()})
		return
	}

	balance, err := h.userSvc.Balance(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, errmsg.UserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": errmsg.UserNotFound.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to retrieve user balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (h *Handler) GetFilteredUserSMS(c *gin.Context) {
	var req dto.GetFilteredUserSMSReq
	req.UserID = c.Param("id")
	if err := c.ShouldBindQuery(&req.Filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg.InvalidRequest.Error()})
		return
	}

	if err := h.validator.FilteredUserSMS(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Filter.Page == 0 {
		req.Filter.Page = 1
	}

	if req.Filter.PageSize == 0 {
		req.Filter.PageSize = 20
	}

	resp, err := h.userSvc.GetFilteredUserSMS(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "failed to retrieve user sms report"})
		return
	}
	if resp.Count < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": errmsg.SMSNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
