package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	ValidationError struct {
		Field string
		Err   error
	}
	ValidationErrors []ValidationError

	// Общий интерфейс для всех валидаторов.
	ValueValidator interface {
		ValidateValue(value reflect.Value) error
	}
	// Интерфейс для валидаторов, принимающих String.
	StringValueValidatorInterface interface {
		ValidateValueString(value string) error
	}
	// Интерфейс для валидаторов, принимающих Int.
	IntValueValidatorInterface interface {
		ValidateValueInt(value int64) error
	}

	// Обобщающая структура для валидаторов string, объявляющая что может реализовывать StringValueValidatorInterface.
	StringValueValidator struct{}
	// Обобщающая структура для валидаторов int.
	IntValueValidator struct{}

	// Собственно валидаторы, реализаторы интерфейсов.
	StringLengthValidator struct {
		length int
		StringValueValidator
	}
	StringRegexpValidator struct {
		expression *regexp.Regexp
		StringValueValidator
	}
	StringInValidator struct {
		in map[string]struct{}
		StringValueValidator
	}
	IntMinValidator struct {
		min int
		IntValueValidator
	}
	IntMaxValidator struct {
		max int
		IntValueValidator
	}
	IntInValidator struct {
		in map[int64]struct{}
		IntValueValidator
	}
)

func (v ValidationErrors) Error() string {
	var buffer bytes.Buffer

	if len(v) > 0 {
		buffer.WriteString("Errors:\n")
	}
	for _, err := range v {
		buffer.WriteString(fmt.Sprintf("- %s\n", err.Err.Error()))
	}
	return buffer.String()
}

// ----------------------------------- Общие ошибки.
type ErrorUnsupportedType struct{ Type string }

func (e ErrorUnsupportedType) Error() string {
	return fmt.Sprintf("value has unsupported type: %s", e.Type)
}

type ErrorUnsupportedValidatorName struct{ Name string }

func (e ErrorUnsupportedValidatorName) Error() string {
	return fmt.Sprintf("validator name '%s' is unsupported", e.Name)
}

type ErrorIncorrectValidator struct {
	Field string
	Err   error
}

func (e ErrorIncorrectValidator) Error() string {
	return fmt.Sprintf("incorrect validator in field '%s': %s", e.Field, e.Err.Error())
}

// -----------------------------------  Ошибки валидации.
type ErrorStringLength struct {
	Length int
	Value  string
}

func (e ErrorStringLength) Error() string {
	return fmt.Sprintf("string value length mismatch %v: %s", e.Length, e.Value)
}

type ErrorMismatchRegular struct {
	Regular string
	Value   string
}

func (e ErrorMismatchRegular) Error() string {
	return fmt.Sprintf("value does not match regular expression %s: %s", e.Regular, e.Value)
}

type ErrorNotInSetValue struct{ Value interface{} }

func (e ErrorNotInSetValue) Error() string {
	return fmt.Sprintf("value %s is not in the set", e.Value)
}

type ErrorNotMinValue struct {
	Min   int
	Value int
}

func (e ErrorNotMinValue) Error() string {
	return fmt.Sprintf("value %v is greater than %v", e.Value, e.Min)
}

type ErrorNotMaxValue struct {
	Max   int
	Value int
}

func (e ErrorNotMaxValue) Error() string {
	return fmt.Sprintf("value %v is less than %v", e.Value, e.Max)
}

// ----------------------------------- Общие переменные.
var (
	ErrorNotStruct                = errors.New("value type is not struct")
	ErrorIncorrectValidatorFormat = errors.New("validator has incorrect format")

	// Индекс валидаторов по типу поля и имени валидатора.
	funcIndexCreateValidator = map[reflect.Kind]map[string]func(cond string) (ValueValidator, error){
		reflect.String: {
			"len":    createStringLengthValidator,
			"regexp": createStringRegexpValidator,
			"in":     createStringInValidator,
		},
		reflect.Int: {
			"min": createIntMinValidator,
			"max": createIntMaxValidator,
			"in":  createIntInValidator,
		},
	}
)

// -----------------------------------  Стартовая точка.
func Validate(v interface{}) error {
	structV := reflect.ValueOf(v)
	// Проверка на то, что входной interface{} - структура.
	if structV.Type().Kind() != reflect.Struct {
		return ErrorNotStruct
	}

	// Проходимся по стартовой структуре structV (и всем вложенным), конвертируем тэги в воркеры-валидаторы validators
	validators := make(map[string][]ValueValidator, 0)
	err := createValidators(structV, validators, "")
	if err != nil {
		return err
	}

	// Проходимся по созданным воркером и выполянем валидацию
	validationErrors := make(ValidationErrors, 0)
	for name, fieldValidators := range validators {
		// Для составных структур по имени валидатора complexName найдем рефлексию внутреннего поля subName
		complexName := strings.Split(name, ".")
		structByName, subName := getInnerStructByName(structV, complexName, 0)
		// structByName - самая внутренняя структура

		// Параметры внутреннего поля для валидации
		field, _ := structByName.Type().FieldByName(subName)
		value := structByName.FieldByName(subName)

		// Несколько проверок для одного и того же поля
		for _, validator := range fieldValidators {
			kind := field.Type.Kind()
			if kind == reflect.Slice || kind == reflect.Array {
				for i := 0; i < value.Len(); i++ {
					valueInner := value.Index(0)
					err = validator.ValidateValue(valueInner)
					if err != nil {
						validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
					}
				}
			} else {
				err = validator.ValidateValue(value)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// Рекурсивно найдем по составному имени самую внутреннюю структуру, и имя внутреннего поля.
func getInnerStructByName(structV reflect.Value, complexName []string, nameLvl int) (reflect.Value, string) {
	subName := complexName[nameLvl]
	structVInner := structV.FieldByName(subName)
	if structVInner.Type().Kind() == reflect.Struct {
		return getInnerStructByName(structVInner, complexName, nameLvl+1)
	}
	return structV, subName
}

// Рекурсивно соберем тэги со структуры, создадим по ним валидаторы и добавим их в общий список воркеров validators.
// extPrefix - путь от самой внешней структуры для текущей.
func createValidators(rv reflect.Value, validators map[string][]ValueValidator, extPrefix string) error {
	for i := 0; i < rv.NumField(); i++ {
		rvField := rv.Field(i)
		structField := rv.Type().Field(i)
		kind := structField.Type.Kind()

		if kind == reflect.Struct {
			tags, _ := structField.Tag.Lookup("validate")
			if tags == "nested" {
				// Создаем валидаторы для внутренних структур
				structPefix := joinSafely(structField.Name, extPrefix)
				err := createValidators(rvField, validators, structPefix)
				if err != nil {
					return err
				}
				continue
			}
		}

		if kind == reflect.Array || kind == reflect.Slice {
			if rvField.Len() == 0 {
				continue
			}
			kind = rvField.Index(0).Type().Kind()
		}
		// создадим валидаторы и добавим их в общий список воркеров validators
		err := addFieldValidators(structField, kind, validators, extPrefix)
		if err != nil {
			return err
		}
	}
	return nil
}

// Соединяем входящие строке только если обе не пустые. Иначе возвращаем без префикса.
func joinSafely(subName string, extPrefix string) string {
	result := subName
	if extPrefix != "" {
		result = strings.Join([]string{extPrefix, result}, ".")
	}
	return result
}

// обработка тега структуры `validate` и создание по указанному правилу валидатора-воркера.
func addFieldValidators(
	structField reflect.StructField,
	kind reflect.Kind,
	validators map[string][]ValueValidator,
	extPrefix string,
) error {
	tags, ok := structField.Tag.Lookup("validate")
	if !ok {
		// Тег для валидации отсутствует. Поэтому проверка не нужна, новых валидаторов не добавляем
		return nil
	}

	// Набор валидаторов, поддержанный для поля типа kind
	funcIndex, ok := funcIndexCreateValidator[kind]
	if !ok {
		return ErrorUnsupportedType{kind.String()}
	}
	fieldValidators := make([]ValueValidator, 0)
	for _, condition := range strings.Split(tags, "|") {
		conditionLexems := strings.Split(condition, ":")

		if len(conditionLexems) != 2 {
			return ErrorIncorrectValidatorFormat
		}
		validatorName := strings.TrimSpace(conditionLexems[0])
		validatorOperand := strings.TrimSpace(conditionLexems[1])

		createValidatorFunc, ok := funcIndex[validatorName]
		if !ok {
			return ErrorUnsupportedValidatorName{validatorName}
		}
		validator, err := createValidatorFunc(validatorOperand)
		if err != nil {
			return ErrorIncorrectValidator{Field: validatorName, Err: err}
		}
		fieldValidators = append(fieldValidators, validator)
	}

	validatorKey := joinSafely(structField.Name, extPrefix)
	validators[validatorKey] = fieldValidators
	return nil
}

// --------------------------------Создание рабочих валидаторов:

func createStringLengthValidator(operand string) (ValueValidator, error) {
	stringLength, err := strconv.Atoi(operand)
	if err != nil {
		return nil, err
	}
	return &StringLengthValidator{length: stringLength}, nil
}

func createStringRegexpValidator(expression string) (ValueValidator, error) {
	reg, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return &StringRegexpValidator{expression: reg}, nil
}

func createStringInValidator(values string) (ValueValidator, error) {
	inValues := make(map[string]struct{})
	for _, val := range strings.Split(values, ",") {
		inValues[val] = struct{}{}
	}
	if len(inValues) == 0 {
		return nil, ErrorIncorrectValidatorFormat
	}
	return &StringInValidator{in: inValues}, nil
}

func createIntMinValidator(operand string) (ValueValidator, error) {
	minInt, err := strconv.Atoi(operand)
	if err != nil {
		return nil, err
	}
	return &IntMinValidator{min: minInt}, nil
}

func createIntMaxValidator(operand string) (ValueValidator, error) {
	maxInt, err := strconv.Atoi(operand)
	if err != nil {
		return nil, err
	}
	return &IntMaxValidator{max: maxInt}, nil
}

func createIntInValidator(cond string) (ValueValidator, error) {
	inValues := make(map[int64]struct{})
	for _, value := range strings.Split(cond, ",") {
		val, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		inValues[int64(val)] = struct{}{}
	}
	return &IntInValidator{in: inValues}, nil
}

// ----------------------------------- Логика валидации:

func (validator *StringLengthValidator) ValidateValue(value reflect.Value) error {
	// Так смело преобразовывать reflect.Value в string нам позволяет вложенная структура StringValueValidator
	s := value.String()
	if len(s) != validator.length {
		return ErrorStringLength{Length: validator.length, Value: s}
	}
	return nil
}

func (validator *StringRegexpValidator) ValidateValue(value reflect.Value) error {
	s := value.String()
	if !validator.expression.MatchString(s) {
		return ErrorMismatchRegular{Regular: validator.expression.String(), Value: s}
	}
	return nil
}

func (validator *StringInValidator) ValidateValue(value reflect.Value) error {
	s := value.String()
	if _, ok := validator.in[s]; !ok {
		return ErrorNotInSetValue{s}
	}
	return nil
}

func (validator IntMinValidator) ValidateValue(value reflect.Value) error {
	// Так смело преобразовывать reflect.Value в int нам позволяет вложенная структура IntValueValidator
	i := value.Int()
	if i < int64(validator.min) {
		return ErrorNotMinValue{Min: validator.min, Value: int(i)}
	}
	return nil
}

func (validator *IntMaxValidator) ValidateValue(value reflect.Value) error {
	i := value.Int()
	if i > int64(validator.max) {
		return ErrorNotMaxValue{Max: validator.max, Value: int(i)}
	}
	return nil
}

func (validator *IntInValidator) ValidateValue(value reflect.Value) error {
	i := value.Int()
	if _, ok := validator.in[i]; !ok {
		return ErrorNotInSetValue{Value: int(i)}
	}
	return nil
}

// Функции-маркеры того, что StringValueValidator может реализовывать StringValueValidatorInterface
// А значит в чистом виде reflect.Value может спокойно преобразовываться в string
// На практике в чистом виде StringValueValidator не должен использоваться, поэтому функция не должна будет паниковать.
func (validator *StringValueValidator) ValidateValueString(_ string) error {
	panic("Must be implemented")
}

// Аналогично для int.
func (validator *IntValueValidator) ValidateValueInt(_ int64) error {
	panic("Must be implemented")
}
