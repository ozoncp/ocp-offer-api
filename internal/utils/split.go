package utils

import (
	"errors"
	"math"
)

// Разделение слайса на батчи - чанки одинкового размера (кроме последнего)
//
// "source" - исходный слайс;
// "batchSize" - количество частей на которые нужно разбить слайс.
func Split(source []int, batchSize int) ([][]int, error) {
	// Проверка на то, что количество батчей больше нуля.
	if batchSize <= 0 {
		return nil, errors.New("the batch size must not be less than zero or equal to zero")
	}

	// Проверка на то, что количество батчей не длиннее размера слайса.
	if batchSize > len(source) {
		return nil, errors.New("the batch size is larger than the slice length")
	}

	// Слайс это структура, которая имеет указатель на выделенный участок памяти с массивом, его длину и вместимость.
	// При передаче в функции слайс копируется, но вместе с указателем на тот же участок памяти,
	// из-за этого при изменении данных возвращаемого слайса меняется данные и в исходном.
	// Исходный код слайса https://golang.org/src/runtime/slice.go#L11
	slice := make([]int, len(source))
	copy(slice, source)

	var result [][]int

	length := len(slice)

	// Количество шагов (батчей), округляем в большую сторону
	step := int(math.Ceil(float64(length) / float64(batchSize)))

	// Разбиваем слайс на части и добавляем в результат
	for i := (0); i < length; i += step {
		j := i + step
		if j > length {
			j = length
		}
		result = append(result, slice[i:j])
	}

	return result, nil
}
