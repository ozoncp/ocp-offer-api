package models_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestOfferString(t *testing.T) {
	t.Parallel()
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name    string       // Название теста
		offer   models.Offer // Исходный слайс
		result  string       // Результат разбивки
		isError bool         // Если должна вернуться ошибка
	}{
		{
			name:   "Valid offer to string",
			offer:  models.Offer{ID: 0, UserID: 1, Grade: 2, TeamID: 3},
			result: "Id: 0, UserId: 1, Grade: 2, TeamId: 3",
		},
		{
			name:   "Empty (default value)",
			offer:  models.Offer{},
			result: "Id: 0, UserId: 0, Grade: 0, TeamId: 0",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := tc.offer.String()
			assert.Equal(t, tc.result, result)
		})
	}
}
