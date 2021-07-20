package utils_test

import (
	"testing"

	"github.com/ozoncp/ocp-offer-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

// Тестирование функции смены ключ-значение на значение-ключ (обратный индекс)
func TestFlipKeyValue(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name    string         // Название теста
		source  map[string]int // Исходный слайс
		result  map[int]string // Результат смены ключ-значение на значение-ключ
		isValid bool           // Является ли тест кейс валидным (иначе проверять на неравенство с результатом)
	}{
		{
			name:    "Empty source & result",
			source:  map[string]int{},
			result:  map[int]string{},
			isValid: true,
		},
		{
			name: "Equal keys",
			source: map[string]int{
				"key-1": 1,
				"key-2": 2,
				"key-3": 3,
			},
			result: map[int]string{
				1: "key-1",
				2: "key-2",
				3: "key-3",
			},
			isValid: true,
		},
		{
			name: "Equal names",
			source: map[string]int{
				"Владислав": 2,
				"Моисей":    1,
				"Роман":     4,
				"Вениамин":  3,
				"Ефим":      0,
				"Геннадий":  5,
			},
			result: map[int]string{
				2: "Владислав",
				1: "Моисей",
				4: "Роман",
				3: "Вениамин",
				0: "Ефим",
				5: "Геннадий",
			},
			isValid: true,
		},
		{
			name: "Empty result",
			source: map[string]int{
				"Дмитрий":  0,
				"Борис":    1,
				"Антон":    2,
				"Владилен": 3,
			},
			result:  map[int]string{},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.FlipKeyValue(tc.source)

			if tc.isValid {
				assert.Equal(t, tc.result, result)
			} else {
				assert.NotEqual(t, tc.result, result)
			}
		})
	}

}
