package order

import (
	"arcs/internal/clients/nats/producer"
	"arcs/internal/configs"
	"arcs/internal/dto"
	"arcs/internal/lock"
	"arcs/internal/models"
	"arcs/internal/repository/order"
	"arcs/internal/repository/sms"
	userSvc "arcs/internal/service/user"
	consts "arcs/internal/utils/const"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

type Svc struct {
	cfg       configs.Config
	userSvc   *userSvc.Svc
	orderRepo *order.Repository
	smsRepo   *sms.Repository
	producer  *producer.Producer
	lock      *lock.Lock
}

func NewOrderSvc(
	cfg configs.Config,
	userSvc *userSvc.Svc,
	orderRepo *order.Repository,
	smsRepo *sms.Repository,
	producer *producer.Producer,
	lock *lock.Lock,
) *Svc {
	_ = producer.EnsureStream()
	return &Svc{
		cfg:       cfg,
		userSvc:   userSvc,
		orderRepo: orderRepo,
		smsRepo:   smsRepo,
		producer:  producer,
		lock:      lock,
	}
}

func (s *Svc) RegisterOrder(ctx context.Context, req dto.OrderRequest) error {
	cost := int64(len(req.Destinations) * s.cfg.Order.SMSCost)

	//lock enough balance to initiation
	if err := s.userSvc.DecreaseBalance(ctx, req.UserID, cost); err != nil {
		return fmt.Errorf("failed to lock sufitiant order cost: %v", err)
	}

	//register order to db
	orderID := uuid.NewString()
	ord := models.Order{
		ID:        orderID,
		CreatedAt: time.Now(),
		UserID:    req.UserID,
		Content:   req.Content,
	}
	if err := s.orderRepo.Submit(ctx, ord); err != nil {
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
			Order:       &ord,
		})
	}

	go s.publish(smsList...)

	return nil
}

func (s *Svc) publish(smss ...models.SMS) {
	var (
		publishedSms []models.SMS
		mu           sync.Mutex
		wg           sync.WaitGroup
	)
	for _, sms := range smss {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pSms := sms.ToProto()
			byteSms, _ := proto.Marshal(pSms)

			if err := s.producer.Publish(s.cfg.Producer.Subjects[0], byteSms, sms.ID); err != nil {
				sms.Status = consts.PendingStatus
				//log.Printf("pending: [%v]", sms.ID)
			} else {
				sms.Status = consts.PublishedStatus
				//log.Printf("published: [%v]", sms.ID)
			}

			mu.Lock()
			publishedSms = append(publishedSms, sms)
			mu.Unlock()
		}()
	}

	//we're waiting for all msgs published asynchronously
	wg.Wait()

	//TODO - replace me with real ctx
	_ = s.smsRepo.CreateSMSBatch(context.Background(), publishedSms)
}

func (s *Svc) RecoverUnPblishSMS() {
	//avoid republish if nats not came up
	if nerr := s.producer.HealthCheck(); nerr != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var initialPoint time.Time

	//check lock
	acquired, err := s.lock.AcquireLock(ctx, consts.RepublishLock)
	if err != nil {
		return
	}

	//check no one other have this lock
	if !acquired {
		return
	}

	//ensure we release lock after crash or end job
	defer s.lock.ReleaseLock(ctx, consts.RepublishLock)

	go func() {
		ticker := time.NewTicker(time.Duration(s.cfg.Basic.RepublishLockDuration) * time.Second / 2)
		for {
			select {
			case <-ticker.C:
				if err := s.lock.ExtendLock(ctx, consts.RepublishLock); err != nil {
					log.Printf("[LOCK] failed to extend lock: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	//s.redisCli.Client.Ping(ctx)
	for {
		//TODO - replace me with real context
		smss, err := s.smsRepo.ListPending(ctx, initialPoint, s.cfg.Basic.PendingProcessBatchSize)
		if err != nil {
			return
		}

		if len(smss) == 0 {
			break
		}

		s.rePublish(smss...)

		initialPoint = smss[len(smss)-1].CreatedAt
	}
}

func (s *Svc) rePublish(smss ...models.SMS) {
	_ = s.producer.EnsureStream()

	var (
		publishedSms []models.SMS
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	for _, sms := range smss {
		wg.Add(1)
		go func() {
			defer wg.Done()

			//TODO - replace me with protobuf
			pSms := sms.ToProto()
			byteSms, _ := proto.Marshal(pSms)

			if err := s.producer.Publish(s.cfg.Producer.Subjects[0], byteSms, sms.ID); err != nil {
				log.Printf("pending: [%v]", sms.ID)
			} else {
				log.Printf("published: [%v]", sms.ID)
				mu.Lock()
				publishedSms = append(publishedSms, sms)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if len(publishedSms) > 0 {
		//TODO - replace me with real context
		_ = s.smsRepo.Update(context.Background(), publishedSms)
	}
}
