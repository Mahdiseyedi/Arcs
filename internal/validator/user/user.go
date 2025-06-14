package user

import (
	"arcs/internal/dto"
	"arcs/internal/utils"
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
