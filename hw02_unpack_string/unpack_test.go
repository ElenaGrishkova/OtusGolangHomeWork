package hw02unpackstring

import (
	"errors"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "a0a0b0", expected: ""},         // Мой новый тест-кейс - все символы зануляются
		{input: "\n3", expected: "\n\n\n"},      // Мой новый тест-кейс - повтор руны с переносом строки
		{input: "a\n0b\n2", expected: "ab\n\n"}, // Мой новый тест-кейс - перенос строки где-то должен убраться, а где-то нет
		{input: ".3", expected: "..."},          // Мой новый тест-кейс - повтор символа не-буквы
		{input: "a", expected: "a"},             // Мой новый тест-кейс - не сломается ли цикл на одном символе
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
