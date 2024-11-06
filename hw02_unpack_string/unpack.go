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
	var isPrevDigitNotEscaped bool
	var isPrevEscaped bool
	for pos, run := range inputRunes {
		// Текущая итерация
		currentSymb := string(run)
		// Строка, которая в этой итерации будет добавлена к общему результату
		// (возможно пустая, возможно задублированный или одиночный символ)
		var addingStr string
		// Символ, из которого будет состоять addingStr повторенный или одиночный
		// (как правило это символ предыдущей итерации, так как информацию про него мы поняли только сейчас)
		toRepeat := prevSymb

		// количество повторений в случае если текущий символ будет неэкранированной цифрой (далее в switch case)
		repeatCount, errCheckDigit := strconv.Atoi(currentSymb)
		// Является ли текущий символ цифрой
		isDigit := errCheckDigit == nil

		// Является ли текущий символ (цифра или слеш) экранированным.
		// Если предыдущий символ - уже являлся экранированным слешем, то на текущий символ не влияет, отчет ведется снова
		isEscaped := prevSymb == `\` && !isPrevEscaped
		// Является ли текущая цифра триггером к повторению предыдущего символа (является ли цифра неэкранированной)
		isDigitForPrevRepeat := isDigit && !isEscaped

		switch {
		case isDigitForPrevRepeat:
			// Текущий символ - неэкранированная цифра обозначающая количество повторений.
			// Проверим что он не первый символ и не вторая цифра
			err2 := checkDigit(pos, isPrevDigitNotEscaped)
			if err2 != nil {
				return "", err2
			}
		case prevSymb == `\` && !isDigit && currentSymb != `\`:
			// Ошибка - экранировать можно только цифры и слэш
			return "", ErrInvalidString
		case isPrevDigitNotEscaped || isEscaped:
			// Если предыдущий символ - цифра или слэш, не записываем этот предыдущий символ в addingStr вообще, повторяем 0 раз
			repeatCount = 0
		default:
			// Текущий символ - не неэкранированная цифра, поэтому предыдущий символ не повторяем.
			// Текущий символ обработаем в следующей итерации
			repeatCount = 1
		}
		// Нужно повторить предыдущий символ repeatCount раз
		addingStr = strings.Repeat(toRepeat, repeatCount)

		result.WriteString(addingStr)

		// Если конец цикла, то отдельным образом приклеим последний символ, если он не неэкранироанная цифра
		if pos == len(inputRunes)-1 && !isDigitForPrevRepeat {
			result.WriteString(currentSymb)
		}

		// Переопределяем переменные для следующей итерации цикла
		prevSymb = currentSymb
		isPrevEscaped = isEscaped
		if isDigitForPrevRepeat {
			isPrevDigitNotEscaped = true
		} else {
			isPrevDigitNotEscaped = false
		}
	}
	return result.String(), nil
}

func checkDigit(pos int, isPrevDigit bool) error {
	if pos == 0 {
		// Строка не может начинаться с числа
		return ErrInvalidString
	}
	if isPrevDigit {
		// Две цифры подряд запрещены
		return ErrInvalidString
	}
	return nil
}
