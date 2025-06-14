package utils

import (
	"arcs/internal/utils/errmsg"
	"github.com/google/uuid"
)

func IsValidUUID(id interface{}) bool {
	idStr, ok := id.(string)
	if !ok {
		return false
	}

	if err := uuid.Validate(idStr); err != nil {
		return false
	}

	return true
}

func BalanceValidator(balance interface{}) error {
	blc, ok := balance.(int64)
	if !ok {
		return errmsg.InvalidTypeAssertion
	}
	if blc < 0 {
		return errmsg.NegativeBalance
	}
	if blc > 1_000_000_000 {
		return errmsg.InsufficientBalance
	}

	return nil
}
