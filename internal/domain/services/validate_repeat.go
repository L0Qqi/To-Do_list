package services

import (
	"errors"
	"strconv"
	"strings"
)

// Проверка формат правила повторения
func ValidateRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}

	if repeat == "y" {
		return nil
	}

	parts := strings.Split(repeat, " ")
	if len(parts) != 2 || parts[0] != "d" {
		return errors.New("некорректный формат repeat, ожидается формат 'd N'")
	}

	_, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("некорректное значение N в repeat, ожидается положительное число")
	}

	return nil
}
