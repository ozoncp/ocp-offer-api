package utils_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"

	"github.com/stretchr/testify/assert"
)

func TestSplitOffersToBatches(t *testing.T) {
	t.Parallel()

	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name      string           // Название теста
		source    []models.Offer   // Исходный слайс
		batchSize uint             // Количество частей на которые нужно разбить слайс
		result    [][]models.Offer // Результат разбивки
		isError   bool             // Если должна вернуться ошибка
	}{
		{
			name: "Batch size 2 & len slice 10",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			batchSize: 2,
			result: [][]models.Offer{
				{
					{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
					{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
					{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
					{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
					{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				},
				{
					{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
					{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
					{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
					{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
					{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
				},
			},
			isError: false,
		},
		{
			name: "Batch size 3 & len slice 10",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			batchSize: 3,
			result: [][]models.Offer{
				{
					{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
					{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
					{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
					{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				},
				{
					{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
					{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
					{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
					{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				},
				{
					{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
					{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
				},
			},
			isError: false,
		},
		{
			name: "Batch size 0 & len slice 10",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			batchSize: 0,
			result:    nil,
			isError:   true,
		},
		{
			name: "Batch size 11 & len slice 10",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			batchSize: 11,
			result:    nil,
			isError:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := utils.SplitOffersToBatches(tc.source, tc.batchSize)

			if tc.isError {
				assert.Error(t, err)
			}

			// Проверка на то, что слайс был скопирован и нет ссылок на исходный массив
			if len(result) > 0 {
				tc.source[0].ID = 9876543210                         // Изменяем значение в исходном слайсе
				assert.NotEqual(t, tc.source[0].ID, result[0][0].ID) // Значение не должно совпадать с результатом
			}

			assert.Equal(t, tc.result, result)
		})
	}
}

func TestConvertOffersSliceToMap(t *testing.T) {
	t.Parallel()

	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		isValid bool                    // Является ли тест кейс валидным (иначе проверять на неравенство с результатом)
		isError bool                    // Если должна вернуться ошибка
		name    string                  // Название теста
		source  []models.Offer          // Исходный слайс
		result  map[uint64]models.Offer // Результат разбивки
	}{
		{

			name:    "Empty source & result",
			source:  []models.Offer{},
			result:  map[uint64]models.Offer{},
			isValid: true,
			isError: false,
		},
		{
			name: "Equal",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			result: map[uint64]models.Offer{
				0: {ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				1: {ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				2: {ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				3: {ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				4: {ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				5: {ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				6: {ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				7: {ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				8: {ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				9: {ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			isValid: true,
			isError: false,
		},
		{
			name: "Empty result",
			source: []models.Offer{
				{ID: 0, UserID: 0, Grade: 0, TeamID: 0},
				{ID: 1, UserID: 1, Grade: 1, TeamID: 1},
				{ID: 2, UserID: 2, Grade: 2, TeamID: 2},
				{ID: 3, UserID: 3, Grade: 3, TeamID: 3},
				{ID: 4, UserID: 4, Grade: 4, TeamID: 4},
				{ID: 5, UserID: 5, Grade: 5, TeamID: 5},
				{ID: 6, UserID: 6, Grade: 6, TeamID: 6},
				{ID: 7, UserID: 7, Grade: 7, TeamID: 7},
				{ID: 8, UserID: 8, Grade: 8, TeamID: 8},
				{ID: 9, UserID: 9, Grade: 9, TeamID: 9},
			},
			result:  map[uint64]models.Offer{},
			isValid: false,
			isError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := utils.ConvertOffersSliceToMap(tc.source)

			if tc.isError {
				assert.Error(t, err)
			}

			// Проверка результата
			if tc.isValid {
				assert.Equal(t, tc.result, result)
			} else {
				assert.NotEqual(t, tc.result, result)
			}
		})
	}
}
