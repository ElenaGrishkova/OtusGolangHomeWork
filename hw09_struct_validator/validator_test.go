package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`

		// Добавлены вложенные структуры
		Owner         User     `validate:"nested"`
		ValidResponse Response `validate:"nested"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	// некорректно заданные правила валидации.
	BrokenStructWrongType struct {
		WrongType bool `validate:"len:5"`
	}
	BrokenStructIncorrectValidatorFormat struct {
		IncorrectValidatorFormat string `validate:"len"`
	}
	BrokenStructUnsupportedValidatorName struct {
		UnsupportedValidatorName string `validate:"replace:i"`
	}
	BrokenStructIncorrectValidator struct {
		IncorrectValidator string `validate:"len:anyLen"`
	}
)

func TestValidate(t *testing.T) {
	var validationError ValidationErrors
	var errorUnsupportedType ErrorUnsupportedType
	var errorUnsupportedValidatorName ErrorUnsupportedValidatorName
	var errorIncorrectValidator ErrorIncorrectValidator
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// 0: Все поля ошибочны
			in: User{
				ID:     "userID_error",
				Phones: []string{"12345678", "111111111111111111111111111111111111111111111111111111111111111111111"},
				Email:  "ema@in@com",
				Role:   "chief",
				meta:   json.RawMessage{1, 'A'},
			},
			expectedErr: validationError,
		},
		{
			// 1: Все поля верные
			in: User{
				ID:     "userID_nice_user_lalalala_1234567890",
				Phones: []string{"12345678901", "11111111111"},
				Age:    23,
				Email:  "email@gm.com",
				Role:   "admin",
				meta:   json.RawMessage{1, 'A'},
			},
			expectedErr: nil,
		},
		{
			// 2: Все поля пустые
			in:          User{},
			expectedErr: validationError,
		},
		{
			// 3: Все поля пустые, кроме не подлежащих валидации
			in: User{
				Name: "TestName",
			},
			expectedErr: validationError,
		},
		{
			// 4: Проверка на int IN - ошибка
			in: Response{
				Code: 302,
			},
			expectedErr: validationError,
		},
		{
			// 5: Проверка на int IN - успех
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			// 6: Проверка структуры, у которой вообще нет полей на валидацию
			in: Token{
				Header: []byte("headerTest"),
			},
			expectedErr: nil,
		},
		{
			// 7: Проверка вложенных структур, все успешно
			in: App{
				Version: "1.2.3",
				Owner: User{
					ID:     "userID_nice_user_lalalala_1234567890",
					Phones: []string{"12345678901", "11111111111"},
					Age:    23,
					Email:  "email@gm.com",
					Role:   "admin",
					meta:   json.RawMessage{1, 'A'},
				},
				ValidResponse: Response{
					Code: 200,
				},
			},
			expectedErr: nil,
		},
		{
			// 8: Проверка вложенных структур, есть ошибки
			in: App{
				Version: "1.2.3",
				Owner: User{
					ID:     "userID_error",
					Phones: []string{"12345678", "111111111111111111111111111111111111111111111111111111111111111111111"},
					Email:  "ema@in@com",
					Role:   "chief",
					meta:   json.RawMessage{1, 'A'},
				},
				ValidResponse: Response{
					Code: 301,
				},
			},
			expectedErr: validationError,
		},
		{
			// 9: Переданное значение - не структура
			in:          "Not Structure",
			expectedErr: ErrorNotStruct,
		},
		{
			// 10: Неизвестный тип поля для валидации
			in: BrokenStructWrongType{
				WrongType: false,
			},
			expectedErr: errorUnsupportedType,
		},
		{
			// 10: Неизвестный тип поля для валидации
			in: BrokenStructIncorrectValidatorFormat{
				IncorrectValidatorFormat: "anyString",
			},
			expectedErr: ErrorIncorrectValidatorFormat,
		},
		{
			// 11: Неизвестная функция валидации
			in: BrokenStructUnsupportedValidatorName{
				UnsupportedValidatorName: "anyString",
			},
			expectedErr: errorUnsupportedValidatorName,
		},
		{
			// 12: Ошибка формата уже известного валидатора
			in: BrokenStructIncorrectValidator{
				IncorrectValidator: "anyString",
			},
			expectedErr: errorIncorrectValidator,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result := Validate(tt.in)

			if tt.expectedErr == nil {
				require.Nil(t, result)
			} else {
				require.ErrorAs(t, result, &tt.expectedErr)
			}
		})
	}
}
