package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

// Знаки препинания, которые следует игнорировать по краям слова.
const punctSymbols = ".,-?!;:()\"'`"

func Top10(input string) []string {
	if input == "" {
		return []string{}
	}
	index := make(map[string]int)
	// Разделим входящую строку на слова
	inputSlice := strings.Fields(input)

	// Вычислим количество повторяющихся слов
	for _, word := range inputSlice {
		pureLexeme := getPureLexeme(word)
		if key, needIndex := getIndexWord(pureLexeme, word); needIndex {
			index[key]++
		}
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

// Преобразовывает изначальное слово - в чистое, пригодное для подсчета.
// Убирает знаки препинания по краям и преобразовывает в нижний регистр.
func getPureLexeme(rawWord string) string {
	pureLexeme := strings.Trim(rawWord, punctSymbols)
	pureLexeme = strings.ToLower(pureLexeme)
	return pureLexeme
}

// Если после сжатия wordForIndex - пустая строка, значит там были только знаки препинания.
// Рассмотрим особые случаи для таких слов.
func getIndexWord(pureLexeme string, rawWord string) (string, bool) {
	if pureLexeme != "" {
		return pureLexeme, true
	}

	if rawWord == "-" {
		// Особый случай: "-" словом не является
		return "", false
	}
	// Остальные случаи - являются словом, пока не поступит отдельного ТЗ.
	return rawWord, true
}
