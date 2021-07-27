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
		isError bool     // Если должна вернуться ошибка
	}{
		{
			name:    "Read one file",
			paths:   []string{"./test/test-file-1.txt"},
			result:  []string{"Test 1"},
			isError: false,
		},
		{
			name:    "Read two file",
			paths:   []string{"./test/test-file-1.txt", "./test/test-file-2.txt"},
			result:  []string{"Test 1", "Test 2"},
			isError: false,
		},
		{
			name:    "Empty",
			paths:   []string{},
			result:  []string{},
			isError: false,
		},
		{
			name:    "Non-existent path",
			paths:   []string{"./test/no-file.txt"},
			result:  nil,
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := main.ReadFiles(tc.paths)

			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.result, result)
		})
	}
}
