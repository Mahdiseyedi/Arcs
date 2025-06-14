package order

import (
	"arcs/internal/clients/nats"
	"arcs/internal/configs"
	"arcs/internal/dto"
	"arcs/internal/models"
	"arcs/internal/repository/order"
	"arcs/internal/repository/sms"
	userSvc "arcs/internal/service/user"
	consts "arcs/internal/utils/const"
	"arcs/internal/utils/errmsg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type Svc struct {
	cfg        configs.Config
	userSvc    *userSvc.Svc
	orderRepo  *order.Repository
	smsRepo    *sms.Repository
	natsClient *nats.Client
}

func NewOrderSvc(
	cfg configs.Config,
	userSvc *userSvc.Svc,
	orderRepo *order.Repository,
	smsRepo *sms.Repository,
	natsClient *nats.Client,
) *Svc {
	return &Svc{
		cfg:        cfg,
		userSvc:    userSvc,
		orderRepo:  orderRepo,
		smsRepo:    smsRepo,
		natsClient: natsClient,
	}
}

func (s *Svc) RegisterOrder(ctx context.Context, req dto.OrderRequest) error {
	//check user exists or not
	balance, err := s.userSvc.Balance(ctx, req.UserID)
	if err != nil {
		return err
	}

	cost := int64(len(req.Destinations) * s.cfg.Order.SMSCost)
	//check for enough balance
	if cost > balance {
		return errmsg.InsufficientBalance
	}

	//lock enough balance to initiation
	if err := s.userSvc.DecreaseBalance(ctx, req.UserID, cost); err != nil {
		return fmt.Errorf("failed to lock sufitiant order cost: %v", err)
	}

	//register order to db
	orderID := uuid.NewString()
	if err := s.orderRepo.Submit(ctx, models.Order{
		ID:        orderID,
		CreatedAt: time.Now(),
		UserID:    req.UserID,
		Content:   req.Content,
	}); err != nil {
		//refund balance for failure
		_ = s.userSvc.ChargeUser(ctx, dto.ChargeUserBalance{
			UserId: req.UserID,
			Amount: cost,
		})
		return fmt.Errorf("failed to submit order: %v", err)
	}

	var smsList []models.SMS
	for _, dest := range req.Destinations {
		smsList = append(smsList, models.SMS{
			ID:          uuid.NewString(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			OrderID:     orderID,
			Destination: dest,
		})
	}

	s.publish(ctx, smsList...)

	return nil
}

func (s *Svc) publish(ctx context.Context, smss ...models.SMS) {
	_ = s.natsClient.EnsureStream()

	var publishedSms []models.SMS
	for _, sms := range smss {
		//TODO - replace me with protobuf
		byteSms, _ := json.Marshal(sms)

		if err := s.natsClient.Publish(s.cfg.Nats.Subjects[0], byteSms); err != nil {
			log.Printf("failed to register job for dst: [%v]", sms)
			sms.Status = consts.PendingStatus
		} else {
			log.Printf("job registered for dst: [%v]", sms)
			sms.Status = consts.PublishedStatus
		}

		publishedSms = append(publishedSms, sms)
	}

	_ = s.smsRepo.CreateSMSBatch(ctx, publishedSms)

}

func (s *Svc) RecoverUnPblishSMS(ctx context.Context) error {
	smss, err := s.smsRepo.ListPending(ctx)
	if err != nil {
		return err
	}

	s.rePublish(ctx, smss...)

	return nil
}

func (s *Svc) rePublish(ctx context.Context, smss ...models.SMS) {
	_ = s.natsClient.EnsureStream()

	var publishedSms []models.SMS

	for _, sms := range smss {
		//TODO - replace me with protobuf
		byteSms, _ := json.Marshal(sms)

		if err := s.natsClient.Publish(s.cfg.Nats.Subjects[0], byteSms); err != nil {
			log.Printf("failed to register job for dst: [%v]", sms)
		} else {
			log.Printf("job registered for dst: [%v]", sms)
			publishedSms = append(publishedSms, sms)
		}
	}

	if len(publishedSms) > 0 {
		_ = s.smsRepo.Update(ctx, publishedSms)
	}
}
