package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	// Регулярка заменена на обычную проверку суффикса
	suffix := "." + domain

	// Полная загрузка файла в память заменена на построчечное чтение
	scanner := bufio.NewScanner(r)
	// Убран массив users, теперь анализируем каждую строку налету, чтобы не держать в памяти все остальные строки
	for scanner.Scan() {
		line := scanner.Bytes()
		user := User{}
		// Замена библиотеки парсинга JSON на easyjson
		if err := user.UnmarshalJSON(line); err != nil {
			// Пропускаем некорректные строки
			continue
		}

		if strings.HasSuffix(user.Email, suffix) {
			emailParts := strings.SplitN(user.Email, "@", 2)
			if len(emailParts) < 2 {
				continue
			}

			key := strings.ToLower(emailParts[1])
			result[key]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	return result, nil
}
