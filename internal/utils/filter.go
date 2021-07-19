package utils

// Фильтрует входной слайс по критерию отсутствия элемента в захардкоженном списке
func Filter(source []int) []int {
	result := make([]int, 0, len(source))
	list := []int{0, 1, 2, 3, 4, 5}

	// Добавляем в результат только те значения, которые отсутствуют в "list"
	for _, value := range source {
		if !Include(list, value) {
			result = append(result, value)
		}
	}

	return result
}

// Поиск элемента в слайсе
func Include(slice []int, search int) bool {
	for _, value := range slice {
		if value == search {
			return true
		}
	}
	return false
}
