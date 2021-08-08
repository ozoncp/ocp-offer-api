package saver

import (
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

	go func() {
		for {
			select {
			case offer := <-s.offersChan:
				s.offers = append(s.offers, offer)

			// Слушаем событие тикера на сохранение в хранилище
			case <-s.notifier.Notify():
				s.flushOffers()

			case <-s.end:
				close(s.offersChan)
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

	if uint(len(s.offers)) >= s.capacity {
		return ErrorMaximumCapacityReached
	}

	s.offersChan <- offer

	return nil
}

func (s *saver) flushOffers() {
	// возращает не добавленные в хранилище элементы и ошибку
	offers, _ := s.flusher.Flush(s.offers)

	// оставляем только не добавленные в хранилище элементы
	s.offers = s.offers[:copy(s.offers, offers)]
}

func (s *saver) Close() {
	s.end <- struct{}{}

	if s.isInitialized {
		s.notifier.Close()
		s.flushOffers()
	}
}
