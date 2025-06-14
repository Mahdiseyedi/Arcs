package order

import (
	"arcs/internal/dto"
	"arcs/internal/utils"
	"arcs/internal/utils/errmsg"
)

type Validator struct {
}

func newOrderValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(req dto.OrderRequest) error {
	if !utils.IsValidUUID(req.UserID) {
		return errmsg.InvalidID
	}
	//TODO - adding validation for content and destinations
	return nil
}
