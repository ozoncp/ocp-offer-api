package utils

import (
	"errors"

	"github.com/ozoncp/ocp-offer-api/internal/models"
)

// SplitOffersToBatches Батчевое разделение слайса на слайс слайсов
//
// "source" - исходный слайс;
// "batchSize" - количество частей на которые нужно разбить слайс.
func SplitOffersToBatches(source []models.Offer, batchSize uint) ([][]models.Offer, error) {
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
	slice := make([]models.Offer, len(source))
	copy(slice, source)

	var result [][]models.Offer

	length := len(slice)

	// Количество шагов (батчей), округляем в большую сторону
	step := int(len(slice)) / int(batchSize)
	if len(slice)%int(batchSize) != 0 {
		step += 1
	}

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

// ConvertOffersSliceToMap - Конвертации слайса от структуры в отображение,
// где ключ идентификатор структуры, а значение сама структура
//
// "source" - исходный слайс;
func ConvertOffersSliceToMap(source []models.Offer) (map[uint64]models.Offer, error) {
	if source == nil {
		return nil, errors.New("source cannot be `nil`")
	}

	result := make(map[uint64]models.Offer, len(source))

	for _, offer := range source {
		result[offer.Id] = offer
	}

	return result, nil
}
