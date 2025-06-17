package buffer

import (
	"arcs/internal/configs"
	"arcs/internal/models"
	"arcs/internal/repository/sms"
	"context"
	"log"
	"sync"
	"time"
)

type StatusFlusher struct {
	cfg           configs.Config
	updates       chan models.StatusUpdate
	wg            sync.WaitGroup
	flushInterval time.Duration
	repo          *sms.Repository
	shutdown      chan struct{}
}

func NewStatusFlusher(
	cfg configs.Config,
	repo *sms.Repository,
) *StatusFlusher {
	sf := &StatusFlusher{
		updates:       make(chan models.StatusUpdate, 80000),
		flushInterval: time.Duration(cfg.Delivery.BufferFlushInterval) * time.Millisecond,
		repo:          repo,
		shutdown:      make(chan struct{}),
	}
	sf.wg.Add(1)
	go sf.run()

	return sf
}

func (sf *StatusFlusher) run() {
	defer sf.wg.Done()

	ticker := time.NewTicker(sf.flushInterval)
	defer ticker.Stop()

	buffer := make([]models.StatusUpdate, 0, 14000)

	for {
		select {
		case update := <-sf.updates:
			buffer = append(buffer, update)
			if len(buffer) >= 7000 {
				sf.flush(buffer)
				buffer = buffer[:0]
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				sf.flush(buffer)
				buffer = buffer[:0]
			}
		case <-sf.shutdown:
			//log.Printf("flush shutdown %v", len(sf.updates))
			if len(buffer) > 0 {
				sf.flush(buffer)
			}
			return
		}
	}
}

func (sf *StatusFlusher) Add(update models.StatusUpdate) {
	select {
	case sf.updates <- update:
		//log.Printf("added new status update to buffer :%v", update)
	default:
		log.Printf("status update full dropping update..., [%v]\n", update.ID)
	}
}

func (sf *StatusFlusher) Stop() {
	close(sf.shutdown)
	sf.wg.Wait()
}

func (sf *StatusFlusher) flush(buffer []models.StatusUpdate) {
	if len(buffer) == 0 {
		return
	}

	if err := sf.repo.BulkUpdate(context.Background(), buffer); err != nil {
		log.Printf("failed to flush: %v", err)
		return
	}
}
