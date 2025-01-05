package nextDate

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	layout := "20060102"
	var parsedDate time.Time
	var err error

	if date == "" {
		parsedDate = time.Now()
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
		return time.Now().Format(layout), nil

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
		} else if days == 1 {
			return parsedDate.Format(layout), nil
		}

		parsedDate = parsedDate.AddDate(0, 0, days)
		for !parsedDate.After(now) {
			parsedDate = parsedDate.AddDate(0, 0, days)
		}
		return parsedDate.Format(layout), nil
	}
	return "", fmt.Errorf("неподдерживаемое правило повторения: %s", repeat)
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	layout := "20060102"

	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	if nowParam == "" || dateParam == "" || repeatParam == "" {
		http.Error(w, "Отсутствуют необходимые параметры", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(layout, nowParam)
	if err != nil {
		http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, "Ошибка вычисления следующей даты", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
