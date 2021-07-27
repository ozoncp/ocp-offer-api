package utils_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/models"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"

	"github.com/stretchr/testify/assert"
)

func TestSplitToBatches(t *testing.T) {
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
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			batchSize: 2,
			result: [][]models.Offer{
				{
					{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
					{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
					{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
					{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
					{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				},
				{
					{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
					{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
					{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
					{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
					{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
				},
			},
			isError: false,
		},
		{
			name: "Batch size 3 & len slice 10",
			source: []models.Offer{
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			batchSize: 3,
			result: [][]models.Offer{
				{
					{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
					{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
					{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
					{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				},
				{
					{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
					{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
					{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
					{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				},
				{
					{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
					{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
				},
			},
			isError: false,
		},
		{
			name: "Batch size 0 & len slice 10",
			source: []models.Offer{
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			batchSize: 0,
			result:    nil,
			isError:   true,
		},
		{
			name: "Batch size 11 & len slice 10",
			source: []models.Offer{
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			batchSize: 11,
			result:    nil,
			isError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.SplitToBatches(tc.source, tc.batchSize)

			if tc.isError {
				assert.Error(t, err)
			}

			// Проверка на то, что слайс был скопирован и нет ссылок на исходный массив
			if len(result) > 0 {
				tc.source[0].Id = 9876543210                         // Изменяем значение в исходном слайсе
				assert.NotEqual(t, tc.source[0].Id, result[0][0].Id) // Значение не должно совпадать с результатом
			}

			assert.Equal(t, tc.result, result)
		})
	}
}

func TestConvertSliceToMap(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name    string                  // Название теста
		source  []models.Offer          // Исходный слайс
		result  map[uint64]models.Offer // Результат разбивки
		isValid bool                    // Является ли тест кейс валидным (иначе проверять на неравенство с результатом)
		isError bool                    // Если должна вернуться ошибка
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
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			result: map[uint64]models.Offer{
				0: {Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				1: {Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				2: {Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				3: {Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				4: {Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				5: {Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				6: {Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				7: {Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				8: {Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				9: {Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			isValid: true,
			isError: false,
		},
		{
			name: "Empty result",
			source: []models.Offer{
				{Id: 0, UserId: 0, Grade: 0, TeamId: 0},
				{Id: 1, UserId: 1, Grade: 1, TeamId: 1},
				{Id: 2, UserId: 2, Grade: 2, TeamId: 2},
				{Id: 3, UserId: 3, Grade: 3, TeamId: 3},
				{Id: 4, UserId: 4, Grade: 4, TeamId: 4},
				{Id: 5, UserId: 5, Grade: 5, TeamId: 5},
				{Id: 6, UserId: 6, Grade: 6, TeamId: 6},
				{Id: 7, UserId: 7, Grade: 7, TeamId: 7},
				{Id: 8, UserId: 8, Grade: 8, TeamId: 8},
				{Id: 9, UserId: 9, Grade: 9, TeamId: 9},
			},
			result:  map[uint64]models.Offer{},
			isValid: false,
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.ConvertSliceToMap(tc.source)

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
