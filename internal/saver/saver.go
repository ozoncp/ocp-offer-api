package saver

import (
	"sync"
	"time"

	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/models"
)

type Saver interface {
	Save(offer models.Offer) error
	Close()
}

type saver struct {
	capacity   uint
	flusher    flusher.Flusher
	tiker      *time.Ticker
	offersChan chan models.Offer
	offers     []models.Offer
	end        chan struct{}
	mu         *sync.Mutex
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(capacity uint, flusher flusher.Flusher, duration time.Duration) (Saver, error) {
	if capacity <= 0 {
		return nil, ErrorCapacionLessOrEqualZero
	}

	if duration <= 0 {
		return nil, ErrorDurationLessOrEqualZero
	}

	if flusher == nil {
		return nil, ErrorFlusherIsNil
	}

	s := &saver{
		capacity:   capacity,
		flusher:    flusher,
		tiker:      time.NewTicker(duration),
		offers:     make([]models.Offer, capacity),
		offersChan: make(chan models.Offer),
		end:        make(chan struct{}),
		mu:         &sync.Mutex{},
	}

	go func() {
		defer func() {
			s.tiker.Stop()
			close(s.offersChan)
			s.flusher.Flush(s.offers)
		}()

		for {
			select {
			case offer := <-s.offersChan:
				s.mu.Lock()
				s.offers = append(s.offers, offer)
				if uint(len(s.offers)) >= s.capacity {
					s.flushOffers()
				}
				s.mu.Unlock()

			// Слушаем событие тикера на сохранение в хранилище
			case <-s.tiker.C:
				s.mu.Lock()
				s.flushOffers()
				s.mu.Unlock()

			case <-s.end:
				return
			}
		}
	}()

	return s, nil
}

func (s *saver) Save(offer models.Offer) error {
	select {
	case <-s.end:
		return ErrorChanelClosed
	default:
		s.offersChan <- offer
	}

	return nil
}

func (s *saver) flushOffers() {
	// возращает не добавленные в хранилище элементы и ошибку
	unsavedOffers, _ := s.flusher.Flush(s.offers)

	// оставляем только не сохранёные
	s.offers = s.offers[:0]
	s.offers = append(s.offers, unsavedOffers...)
}

func (s *saver) Close() {
	close(s.end)
}
