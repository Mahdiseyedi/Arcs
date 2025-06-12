package order

import (
	"arcs/internal/configs"
	"arcs/internal/dto"
	"arcs/internal/models"
	"arcs/internal/repository/order"
	userSvc "arcs/internal/service/user"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Svc struct {
	cfg       configs.Config
	userSvc   *userSvc.Svc
	orderRepo *order.Repository
}

func NewOrderSvc(
	cfg configs.Config,
	userSvc *userSvc.Svc,
	orderRepo *order.Repository,
) *Svc {
	return &Svc{
		cfg:       cfg,
		userSvc:   userSvc,
		orderRepo: orderRepo,
	}
}

func (s *Svc) RegisterOrder(ctx context.Context, req dto.OrderRequest) error {
	//check user exists or not
	balance, err := s.userSvc.Balance(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("not found any user: %w", err)
	}

	cost := int64(len(req.Destinations) * s.cfg.Order.SMSCost)
	//check for enough balance
	if cost > balance {
		return fmt.Errorf("not enough balance to register order")
	}

	//lock enough balance to initiation
	if err := s.userSvc.DecreaseBalance(ctx, req.UserID, cost); err != nil {
		return fmt.Errorf("failed to lock sufitiant order cost: %v", err)
	}

	//register order to db
	if err := s.orderRepo.Submit(ctx, models.Order{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now(),
		UserID:       req.UserID,
		Content:      req.Content,
		Destinations: req.Destinations,
	}); err != nil {
		//refund balance for failer
		_ = s.userSvc.ChargeUser(ctx, req.UserID, cost)
		return fmt.Errorf("failed to submit order: %v", err)
	}

	//TODO - adding job for each dst to job queue
	for _, dest := range req.Destinations {
		log.Printf("job registered for dst: [%v]", dest)
	}

	return nil
}
