package saver

import "errors"

var (
	// ErrorAlreadyInitialized - saver уже инициализирован.
	ErrorAlreadyInitialized = errors.New("saver is already initialized")

	// ErrorNotInitialized - saver не был инициализирован.
	ErrorNotInitialized = errors.New("saver is not initialized")

	// ErrorMaximumCapacityReached - достигнута максимальная емкость слайса в saver.
	ErrorMaximumCapacityReached = errors.New("cannot add new item, slice capacity is equal to item count")

	// ErrorCapacionLessOrEqualZero - емкость не может быть меньше или равна нулю.
	ErrorCapacionLessOrEqualZero = errors.New("capation cannot be less than or equal to zero")

	// ErrorFlusherIsNil - flusher является nil.
	ErrorFlusherIsNil = errors.New("flusher is nil")

	// ErrorDurationLessOrEqualZero - duration не может быть меньше или равен нулю.
	ErrorDurationLessOrEqualZero = errors.New("duration cannot be less than or equal to zero")

	// ErrorChanelClosed - канал закрыт.
	ErrorChanelClosed = errors.New("chanel is closed")
)
