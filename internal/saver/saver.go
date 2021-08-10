package saver

import (
	"sync"

	"github.com/ozoncp/ocp-offer-api/internal/flusher"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/notifier"
)

type Saver interface {
	Save(offer models.Offer) error
	Init() error
	Close()
}

type saver struct {
	capacity      uint
	flusher       flusher.Flusher
	notifier      notifier.Notifier
	offersChan    chan models.Offer
	offers        []models.Offer
	end           chan struct{}
	wait          *sync.WaitGroup
	mu            *sync.Mutex
	isInitialized bool
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(capacity uint, flusher flusher.Flusher, notifier notifier.Notifier) (Saver, error) {
	if capacity <= 0 {
		return nil, ErrorCapacionLessOrEqualZero
	}

	if flusher == nil {
		return nil, ErrorFlusherIsNil
	}

	if notifier == nil {
		return nil, ErrorNotifierIsNil
	}

	return &saver{
		capacity:      capacity,
		flusher:       flusher,
		notifier:      notifier,
		offers:        make([]models.Offer, capacity),
		offersChan:    make(chan models.Offer),
		end:           make(chan struct{}),
		wait:          &sync.WaitGroup{},
		mu:            &sync.Mutex{},
		isInitialized: false,
	}, nil
}

func (s *saver) Init() error {
	if s.isInitialized {
		return ErrorAlreadyInitialized
	}

	err := s.notifier.Init()
	if err != nil {
		return err
	}

	s.offers = s.offers[:0]

	s.wait.Add(1)

	go func() {
		defer func() {
			close(s.offersChan)
			s.notifier.Close()
			s.flusher.Flush(s.offers)
			s.wait.Done()
		}()

		for {

			select {
			case offer := <-s.offersChan:
				s.mu.Lock()
				s.offers = append(s.offers, offer)
				s.mu.Unlock()

			// Слушаем событие тикера на сохранение в хранилище
			case <-s.notifier.Notify():
				s.flushOffers()

			case <-s.end:
				return
			}
		}
	}()

	s.isInitialized = true

	return nil
}

func (s *saver) Save(offer models.Offer) error {
	if !s.isInitialized {
		return ErrorNotInitialized
	}

	select {
	case <-s.end:
		return ErrorChanelClosed
	default:
		s.offersChan <- offer
	}

	return nil
}

func (s *saver) flushOffers() {
	s.mu.Lock()
	// возращает не добавленные в хранилище элементы и ошибку
	offers, _ := s.flusher.Flush(s.offers)

	// оставляем только не добавленные в хранилище элементы
	s.offers = s.offers[:copy(s.offers, offers)]
	s.mu.Unlock()
}

func (s *saver) Close() {
	close(s.end)
	s.wait.Wait()
}
