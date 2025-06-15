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

	go s.publish(smsList...)

	return nil
}

func (s *Svc) publish(smss ...models.SMS) {
	_ = s.natsClient.EnsureStream()
	var publishedSms []models.SMS
	for _, sms := range smss {
		//TODO - replace me with protobuf
		byteSms, _ := json.Marshal(sms)

		if err := s.natsClient.Publish(s.cfg.Nats.Subjects[0], byteSms, sms.ID); err != nil {
			log.Printf("failed to register job for dst: [%v]", sms)
			sms.Status = consts.PendingStatus
		} else {
			log.Printf("job registered for dst: [%v]", sms)
			sms.Status = consts.PublishedStatus
		}

		publishedSms = append(publishedSms, sms)
	}

	//TODO - replace me with real ctx
	_ = s.smsRepo.CreateSMSBatch(context.Background(), publishedSms)
}

func (s *Svc) RecoverUnPblishSMS() error {
	//avoid republish if nats not came up
	if nerr := s.natsClient.HealthCheck(); nerr != nil {
		return nerr
	}

	log.Println("start processing pending sms...")
	ctx := context.Background()
	var initialPoint time.Time

	for {
		//TODO - replace me with real context
		smss, err := s.smsRepo.ListPending(ctx, initialPoint, s.cfg.Basic.PendingProcessBatchSize)
		if err != nil {
			return err
		}

		if len(smss) == 0 {
			break
		}

		s.rePublish(smss...)

		initialPoint = smss[len(smss)-1].CreatedAt
		log.Printf("len processed sms: %v", len(smss))
	}

	log.Println("process pending sms finished...")
	return nil
}

func (s *Svc) rePublish(smss ...models.SMS) {
	_ = s.natsClient.EnsureStream()

	var publishedSms []models.SMS

	for _, sms := range smss {
		//TODO - replace me with protobuf
		byteSms, _ := json.Marshal(sms)

		if err := s.natsClient.Publish(s.cfg.Nats.Subjects[0], byteSms, sms.ID); err != nil {
			log.Printf("failed to register job for dst: [%v]", sms)
		} else {
			log.Printf("job registered for dst: [%v]", sms)
			publishedSms = append(publishedSms, sms)
		}
	}

	if len(publishedSms) > 0 {
		//TODO - replace me with real context
		_ = s.smsRepo.Update(context.Background(), publishedSms)
	}
}
