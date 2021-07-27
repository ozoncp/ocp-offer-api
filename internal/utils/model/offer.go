package utils

import (
	"errors"
	"math"

	"github.com/ozoncp/ocp-offer-api/internal/model"
)

// Батчевое разделение слайса на слайс слайсов
//
// "source" - исходный слайс;
// "batchSize" - количество частей на которые нужно разбить слайс.
func SplitToBulks(source []model.Offer, batchSize uint) ([][]model.Offer, error) {
	// Проверка на то, что количество батчей больше нуля.
	if batchSize <= 0 {
		return nil, errors.New("the batch size must not be less than zero or equal to zero")
	}

	// Проверка на то, что количество батчей не длиннее размера слайса.
	if batchSize > uint(len(source)) {
		return nil, errors.New("the batch size is larger than the slice length")
	}

	// Слайс это структура, которая имеет указатель на выделенный участок памяти с массивом, его длину и вместимость.
	// При передаче в функции слайс копируется, но вместе с указателем на тот же участок памяти,
	// из-за этого при изменении данных возвращаемого слайса меняется данные и в исходном.
	// Исходный код слайса https://golang.org/src/runtime/slice.go#L11
	slice := make([]model.Offer, len(source))
	copy(slice, source)

	var result [][]model.Offer

	length := len(slice)

	// Количество шагов (батчей), округляем в большую сторону
	step := int(math.Ceil(float64(length) / float64(batchSize)))

	// Разбиваем слайс на части и добавляем в результат
	for i := 0; i < length; i += step {
		j := i + step
		if j > length {
			j = length
		}
		result = append(result, slice[i:j])
	}

	return result, nil
}

// Конвертации слайса от структуры в отображение,
// где ключ идентификатор структуры, а значение сама структура
//
// "source" - исходный слайс;
func ConvertSliceToMap(source []model.Offer) (map[uint64]model.Offer, error) {

	result := make(map[uint64]model.Offer, len(source))

	for _, offer := range source {
		result[offer.Id] = offer
	}

	return result, nil
}
