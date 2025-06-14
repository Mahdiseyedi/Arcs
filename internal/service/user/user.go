package user

import (
	"arcs/internal/dto"
	"arcs/internal/repository/balance"
	"arcs/internal/repository/user"
	"context"
	"github.com/google/uuid"
)

type Svc struct {
	userRepo    *user.Repository
	balanceRepo *balance.Repository
}

func NewUserSvc(
	userRepo *user.Repository,
	balance *balance.Repository,
) *Svc {
	return &Svc{
		userRepo:    userRepo,
		balanceRepo: balance,
	}
}

func (s *Svc) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	uid := uuid.New().String()

	if _, err := s.userRepo.CreateUser(ctx, uid, req.Balance); err != nil {
		return err
	}

	if err := s.balanceRepo.Set(ctx, uid, req.Balance); err != nil {
		return err
	}

	return nil
}

func (s *Svc) Balance(ctx context.Context, uid string) (int64, error) {
	return s.balanceRepo.Get(ctx, uid)
}

func (s *Svc) ChargeUser(ctx context.Context, req dto.ChargeUserBalance) error {
	return s.balanceRepo.Increase(ctx, req.UserId, req.Amount)
}

func (s *Svc) DecreaseBalance(ctx context.Context, uid string, amount int64) error {
	return s.balanceRepo.Decrease(ctx, uid, amount)
}
