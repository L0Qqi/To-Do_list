package taskhandlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/L0Qqi/go_final_project/internal/domain/services/nextDate"
)

// HandleNextDate обрабатывает запросы на вычисление следующей даты.
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
		http.Error(w, "Неправильный формат текущей даты", http.StatusBadRequest)
		return
	}

	nextDateStr, err := nextDate.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDateStr)
}
