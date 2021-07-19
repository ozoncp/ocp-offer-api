package utils

// Смена ключ-значение на значение-ключ (обратный индекс)
func FlipKeyValue(source map[string]int) map[int]string {
	result := make(map[int]string, len(source))

	for key, value := range source {
		result[value] = key
	}

	return result
}
