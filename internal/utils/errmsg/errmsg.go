package errmsg

import "fmt"

var (
	InvalidTypeAssertion = fmt.Errorf("invalid type assertion")
	InvalidID            = fmt.Errorf("invalid ID")
	InvalidRequest       = fmt.Errorf("invalid request")
)

var (
	FailedCreateUser    = fmt.Errorf("failed to create user")
	UserNotFound        = fmt.Errorf("user not found")
	NegativeBalance     = fmt.Errorf("negative balance")
	InsufficientBalance = fmt.Errorf("insufficient balance")
	SMSNotFound         = fmt.Errorf("sms not found")
)
