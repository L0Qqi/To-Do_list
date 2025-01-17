package nextDate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDateAdd(now time.Time, date string, repeat string) (string, error) {
	layout := "20060102"
	parsedDate := now

	if date != "" {
		var err error
		parsedDate, err = time.Parse(layout, date)
		if err != nil {
			return "", fmt.Errorf("неправильный формат даты: %v", err)
		}
	}

	// Если repeat пустой, возвращаем текущую или указанную дату
	if repeat == "" {
		if parsedDate.After(now) {
			return parsedDate.Format(layout), nil
		}
		return now.Format(layout), nil
	}

	if repeat == "y" {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
		return parsedDate.Format(layout), nil
	}

	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", fmt.Errorf("неправильный формат правила повторения: %s", repeat)
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неправильное количество дней: %v", err)
		}

		if days == 1 {
			return parsedDate.Format(layout), nil
		}

		for parsedDate.Before(now) {
			parsedDate = parsedDate.AddDate(0, 0, days)
		}
		return parsedDate.Format(layout), nil
	}

	return "", fmt.Errorf("неподдерживаемое правило повторения: %s", repeat)
}
