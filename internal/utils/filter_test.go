package utils_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

// Тестирование функции фильтрации по захардкоженному списку
// При изменении списка в "Filter", тест кейсы могут перестать быть успешными
// Делать переменную доступной из другого пакета и изменять её для тестов может повлечь за собой непредвиденные ошибки в рантайме
func TestFilter(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name    string // Название теста
		source  []int  // Исходный слайс
		result  []int  // Результат фильтрации
		isValid bool   // Является ли тест кейс валидным (иначе проверять на неравенство с результатом)
	}{
		{
			name:    "Empty source & result",
			source:  []int{},
			result:  []int{},
			isValid: true,
		},
		{
			name:    "Equal source & result",
			source:  []int{1, 2, 3, 4},
			result:  []int{1, 2, 3, 4},
			isValid: false,
		},
		{
			name:    "Filtering not include result",
			source:  []int{1, 5, 7, 9},
			result:  []int{7, 9},
			isValid: true,
		},
		{
			name:    "Empty result",
			source:  []int{1, 5, 3, 2},
			result:  []int{},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.Filter(tc.source)

			if tc.isValid {
				assert.Equal(t, tc.result, result)
			} else {
				assert.NotEqual(t, tc.result, result)
			}
		})
	}

}

// Тестирование функции поиска элемента в слайсе
func TestInclude(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name   string // Название теста
		source []int  // Исходный слайс
		search int    // Искомый элемент
		result bool   // Результат поиска
	}{
		{
			name:   "Empty source",
			source: []int{},
			search: 1,
			result: false,
		},
		{
			name:   "Include",
			source: []int{1, 2, 3, 4},
			search: 1,
			result: true,
		},
		{
			name:   "Not include",
			source: []int{1, 2, 3, 4},
			search: 9,
			result: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.Include(tc.source, tc.search)
			assert.Equal(t, tc.result, result)
		})
	}
}
