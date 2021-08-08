package notifier

import "errors"

var (
	// ErrorAlreadyInitialized - notifier уже инициализирован
	ErrorAlreadyInitialized = errors.New("notifier is already initialized")

	// ErrorDurationLessOrEqualZero - продолжительность не может быть меньше или равна нулю
	ErrorDurationLessOrEqualZero = errors.New("duration cannot be less than or equal to zero")
)
