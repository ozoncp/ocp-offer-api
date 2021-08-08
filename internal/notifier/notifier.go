package notifier

import (
	"time"
)

type Notifier interface {
	Notify() <-chan struct{}
	Init() error
	Close()
}

type notifier struct {
	duration      time.Duration
	notify        chan struct{}
	end           chan struct{}
	isInitialized bool
}

func NewNotifier(duration time.Duration) (Notifier, error) {

	if duration <= 0 {
		return nil, ErrorDurationLessOrEqualZero
	}

	return &notifier{
		duration:      duration,
		notify:        make(chan struct{}),
		end:           make(chan struct{}),
		isInitialized: false,
	}, nil
}

// Init - инициализируем
func (n *notifier) Init() error {
	if n.isInitialized {
		return ErrorAlreadyInitialized
	}

	go func() {
		// Создаём тикер и устанавлеваем значение указанное при создании
		ticker := time.NewTicker(n.duration)
		defer ticker.Stop()
		defer close(n.notify)
		defer close(n.end)

		for {
			select {
			case <-ticker.C:
				n.notify <- struct{}{}
			case <-n.end:
				return
			}
		}
	}()

	n.isInitialized = true

	return nil
}

func (n *notifier) Notify() <-chan struct{} {
	return n.notify
}

func (n *notifier) Close() {
	n.end <- struct{}{}
}
