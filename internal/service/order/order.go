package order

import (
	"arcs/internal/clients/nats/producer"
	"arcs/internal/configs"
	"arcs/internal/dto"
	"arcs/internal/models"
	"arcs/internal/repository/order"
	userSvc "arcs/internal/service/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Svc struct {
	cfg       configs.Config
	userSvc   *userSvc.Svc
	orderRepo *order.Repository
	producer  *producer.NatsProducerClient
}

func NewOrderSvc(
	cfg configs.Config,
	userSvc *userSvc.Svc,
	orderRepo *order.Repository,
	producer *producer.NatsProducerClient,
) *Svc {
	if err := producer.EnsureStream(); err != nil {
		log.Fatal(err)
	}
	
	return &Svc{
		cfg:       cfg,
		userSvc:   userSvc,
		orderRepo: orderRepo,
		producer:  producer,
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
	orderID := uuid.NewString()
	if err := s.orderRepo.Submit(ctx, models.Order{
		ID:           orderID,
		CreatedAt:    time.Now(),
		UserID:       req.UserID,
		Content:      req.Content,
		Destinations: req.Destinations,
	}); err != nil {
		//refund balance for failer
		_ = s.userSvc.ChargeUser(ctx, req.UserID, cost)
		return fmt.Errorf("failed to submit order: %v", err)
	}

	for _, dest := range req.Destinations {
		sms := models.SMS{
			//CreatedAt:   time.,
			ID:          uuid.NewString(),
			UserID:      req.UserID,
			OrderID:     orderID,
			Destination: dest,
			//Status:      "",
		}
		//TODO - replace me with protobuf
		byteSms, _ := json.Marshal(sms)

		if err := s.producer.Publish(s.cfg.Nats.Subjects[0], byteSms); err != nil {
			log.Printf("failed to register job for dst: [%v]", dest)
		} else {
			log.Printf("job registered for dst: [%v]", dest)
		}
	}

	return nil
}
