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
	//Если дата не указана, возвращаем текущую
	if date == "" {
		parsedDate = now
	} else {
		parsedDate, err = time.Parse(layout, date)
		if err != nil {
			return "", fmt.Errorf("неправильный формат даты: %v", err)
		}
	}
	//Если правило повторение не указано, возвращаем текущую дату
	if repeat == "" {
		if parsedDate.After(now) {
			return parsedDate.Format(layout), nil
		}
		return now.Format(layout), nil

	}
	//Обработка случая ежегодным повторением
	if repeat == "y" {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
		return parsedDate.Format(layout), nil
	}
	//Обработка случая с повторением каждые n дней
	if strings.HasPrefix(repeat, "d ") {
		repeatDate := strings.Split(repeat, " ")
		if len(repeatDate) != 2 {
			return "", fmt.Errorf("неправильный формат правила повторения: %s", repeat)
		}
		days, err := strconv.Atoi(repeatDate[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неправильное количество дней: %v", err)
		}

		//fmt.Printf("Начальная дата: %s, now: %s, шаг: %d\n", parsedDate.Format(layout), now.Format(layout), days)

		parsedDate = parsedDate.AddDate(0, 0, days)
		for parsedDate.Format("20060102") == now.Format("20060102") || parsedDate.Before(now) {
			//fmt.Printf("Дата до обновления: %s\n", parsedDate.Format(layout))
			parsedDate = parsedDate.AddDate(0, 0, days)
			//fmt.Printf("Дата после обновления: %s\n", parsedDate.Format(layout))
		}

		//fmt.Printf("Конечная дата: %s\n", parsedDate.Format(layout))
		return parsedDate.Format(layout), nil

	}
	return "", fmt.Errorf("неподдерживаемое правило повторения: %s", repeat)
}
