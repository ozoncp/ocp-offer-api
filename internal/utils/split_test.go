package utils_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name      string  // Название теста
		source    []int   // Исходный слайс
		batchSize int     // Количество частей на которые нужно разбить слайс
		result    [][]int // Результат разбивки
		isError   bool    // Если должна вернуться ошибка
	}{
		{
			name:      "Batch size 2 & len slice 10",
			source:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			batchSize: 2,
			result:    [][]int{{0, 1, 2, 3, 4}, {5, 6, 7, 8, 9}},
			isError:   false,
		},
		{
			name:      "Batch size 3 & len slice 10",
			source:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			batchSize: 3,
			result:    [][]int{{0, 1, 2, 3}, {4, 5, 6, 7}, {8, 9}},
			isError:   false,
		},
		{
			name:      "Batch size 0 & len slice 10",
			source:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			batchSize: 0,
			result:    nil,
			isError:   true,
		},
		{
			name:      "Batch size 11 & len slice 10",
			source:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			batchSize: 11,
			result:    nil,
			isError:   true,
		},
		{
			name:      "Batch size -11 & len slice 10",
			source:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			batchSize: -11,
			result:    nil,
			isError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.Split(tc.source, tc.batchSize)

			if tc.isError {
				assert.Error(t, err)
			}

			assert.Equal(t, tc.result, result)
		})
	}
}
