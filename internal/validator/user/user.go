package user

import (
	"arcs/internal/dto"
	"arcs/internal/utils"
	consts "arcs/internal/utils/const"
	"arcs/internal/utils/errmsg"
)

type Validator struct{}

func NewUserValidator() *Validator {
	return &Validator{}
}

func (v *Validator) CreateUser(req dto.CreateUserRequest) error {
	return utils.BalanceValidator(req.Balance)
}

func (v *Validator) ChargeUser(req dto.ChargeUserBalance) error {
	if !utils.IsValidUUID(req.UserId) {
		return errmsg.InvalidID
	}

	if err := utils.BalanceValidator(req.Amount); err != nil {
		return err
	}

	return nil
}

func (v *Validator) FilteredUserSMS(req dto.GetFilteredUserSMSReq) error {
	if !utils.IsValidUUID(req.UserID) {
		return errmsg.InvalidID
	}

	if req.Filter.Page < 0 {
		return errmsg.InvalidPage
	}

	if req.Filter.PageSize <= 0 || req.Filter.PageSize > 100 {
		return errmsg.InvalidPageSize
	}

	validStatuses := map[string]bool{
		consts.PublishedStatus: true,
		consts.PendingStatus:   true,
		consts.FailedStatus:    true,
		consts.DeliveredStatus: true,
	}
	if req.Filter.Status != "" {
		if _, ok := validStatuses[req.Filter.Status]; !ok {
			return errmsg.InvalidStatus
		}
	}

	if req.Filter.StartDate != nil && req.Filter.EndDate != nil {
		if req.Filter.StartDate.After(*req.Filter.EndDate) {
			return errmsg.InvalidRequest
		}
	}

	return nil
}
