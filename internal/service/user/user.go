package user

import (
	"arcs/internal/dto"
	"arcs/internal/repository/sms"
	"arcs/internal/repository/user"
	"context"
	"github.com/google/uuid"
)

type Svc struct {
	userRepo *user.Repository
	smsRepo  *sms.Repository
}

func NewUserSvc(
	userRepo *user.Repository,
	smsRepo *sms.Repository,
) *Svc {
	return &Svc{
		userRepo: userRepo,
		smsRepo:  smsRepo,
	}
}

func (s *Svc) CreateUser(ctx context.Context, req dto.CreateUserRequest) (resp dto.CreateUserResponse, err error) {
	uid := uuid.New().String()

	if err = s.userRepo.CreateUser(ctx, uid, req.Balance); err != nil {
		return
	}

	return dto.CreateUserResponse{UserID: uid}, nil
}

func (s *Svc) Balance(ctx context.Context, uid string) (int64, error) {
	return s.userRepo.GetUserBalance(ctx, uid)
}

func (s *Svc) ChargeUser(ctx context.Context, req dto.ChargeUserBalance) error {
	return s.userRepo.IncreaseBalance(ctx, req.UserId, req.Amount)
}

func (s *Svc) DecreaseBalance(ctx context.Context, uid string, amount int64) error {
	return s.userRepo.DecreaseBalance(ctx, uid, amount)
}

func (s *Svc) GetFilteredUserSMS(ctx context.Context, req dto.GetFilteredUserSMSReq) (dto.GetFilteredUserSMSResp, error) {
	smss, cnt, err := s.smsRepo.GetUserSMS(ctx, req.UserID, req.Filter)
	if err != nil {
		return dto.GetFilteredUserSMSResp{}, err
	}

	return dto.GetFilteredUserSMSResp{
		SMS:   smss,
		Count: cnt,
	}, nil
}
