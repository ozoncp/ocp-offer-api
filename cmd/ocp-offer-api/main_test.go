package main_test

import (
	"testing"

	main "github.com/ozoncp/ocp-offer-api/cmd/ocp-offer-api"
	"github.com/stretchr/testify/assert"
)

func TestReadFiles(t *testing.T) {
	// Проверка нескольких тестовых кейсов
	testCases := []struct {
		name    string   // Название теста
		paths   []string // Пути до файлов
		result  []string // Результат чтения файлов
		isValid bool     // Является ли тест кейс валидным (иначе проверять на неравенство с результатом)
	}{
		{
			name:    "Read one file",
			paths:   []string{"./test/test-file-1.txt"},
			result:  []string{"Test 1"},
			isValid: true,
		},
		{
			name:    "Read two file",
			paths:   []string{"./test/test-file-1.txt", "./test/test-file-2.txt"},
			result:  []string{"Test 1", "Test 2"},
			isValid: true,
		},
		{
			name:    "Empty",
			paths:   []string{},
			result:  []string{},
			isValid: true,
		},
		{
			name:    "Not equal result",
			paths:   []string{"./test/test-file-1.txt"},
			result:  []string{},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := main.ReadFiles(tc.paths)

			if tc.isValid {
				assert.Equal(t, tc.result, result)
			} else {
				assert.NotEqual(t, tc.result, result)
			}
		})
	}
}
