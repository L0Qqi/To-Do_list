package nextDate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	layout := "20060102"
	var parsedDate time.Time
	var err error

	if date == "" {
		parsedDate = now
	} else {
		parsedDate, err = time.Parse(layout, date)
		if err != nil {
			return "", fmt.Errorf("неправильный формат даты: %v", err)
		}
	}
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
		repeatDate := strings.Split(repeat, " ")
		if len(repeatDate) != 2 {
			return "", fmt.Errorf("неправильный формат правила повторения: %s", repeat)
		}
		days, err := strconv.Atoi(repeatDate[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неправильное количество дней: %v", err)
		}
		// if parsedDate.Equal(now) {
		// 	return parsedDate.Format(layout), nil
		// }

		//parsedDate = parsedDate.AddDate(0, 0, days)
		// for !parsedDate.After(now) {
		// 	parsedDate = parsedDate.AddDate(0, 0, days)
		// }
		// return parsedDate.Format(layout), nil

		fmt.Printf("Начальная дата: %s, now: %s, шаг: %d\n", parsedDate.Format(layout), now.Format(layout), days)

		for parsedDate.Before(now) || parsedDate.Equal(now) {
			parsedDate = parsedDate.AddDate(0, 0, days)
			fmt.Printf("Обновлённая дата: %s\n", parsedDate.Format(layout))
		}
		return parsedDate.Format(layout), nil

	}
	return "", fmt.Errorf("неподдерживаемое правило повторения: %s", repeat)
}
