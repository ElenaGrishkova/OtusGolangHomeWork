package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}
	index := make(map[string]int)
	// Разделим входящую строку на слова
	inputSlice := strings.Fields(input)

	// Вычислим количество повторяющихся слов
	for _, word := range inputSlice {
		index[word]++
	}
	fmt.Println("Index=", index)

	// Получим и отсортируем список ключей по убыванию значений в index
	keySlice := make([]string, 0, len(index))
	for key := range index {
		keySlice = append(keySlice, key)
	}
	sort.Slice(keySlice, func(i, j int) bool {
		if index[keySlice[i]] != index[keySlice[j]] {
			// Для слов с различной частотой сортировка по частоте
			return index[keySlice[i]] > index[keySlice[j]]
		}
		// Для слов с одинаковой частотой сортировка лексикографическая
		return keySlice[i] < keySlice[j]
	})

	// Возвращаем первые топ 10 ключей (или меньше)
	keySlice = keySlice[:min(10, len(keySlice))]
	return keySlice
}
