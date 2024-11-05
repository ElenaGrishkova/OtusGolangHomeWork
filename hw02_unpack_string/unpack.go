package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	inputRunes := []rune(input)
	var prevSymb string
	for pos, run := range inputRunes {
		currentSymb := string(run)
		var addingStr string
		repeatCount, errCheckDigit := strconv.Atoi(currentSymb)
		_, errCheckPrevDigit := strconv.Atoi(prevSymb)

		switch {
		case errCheckDigit == nil:
			// Текущий символ - цифра. Проверим что он не первый символ и не вторая цифра
			err2 := checkDigit(pos, prevSymb)
			if err2 != nil {
				return "", err2
			}
		case errCheckPrevDigit == nil:
			// Если предыдущий символ - цифра, не записываем его в addingStr вообще
			repeatCount = 0
		default:
			// Текущий символ - не цифра, поэтому берем предыдущий символ 1 раз. Текущий символ обработаем в следующей итерации
			repeatCount = 1
		}
		// Нужно повторить предыдущий символ repeatCount раз
		addingStr = strings.Repeat(prevSymb, repeatCount)

		result.WriteString(addingStr)
		prevSymb = currentSymb

		// Если конец цикла, то отдельным образом обработаем последний символ, если он не цифра
		if pos == len(inputRunes)-1 && errCheckDigit != nil {
			result.WriteString(currentSymb)
		}
	}
	return result.String(), nil
}

func checkDigit(pos int, prevSymb string) error {
	if pos == 0 {
		// Строка не может начинаться с числа
		return ErrInvalidString
	}
	_, err := strconv.Atoi(prevSymb)
	if err == nil {
		// Две цифры подряд запрещены
		return ErrInvalidString
	}
	return nil
}
