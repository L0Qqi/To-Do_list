package taskhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/L0Qqi/go_final_project/internal/app"
	"github.com/L0Qqi/go_final_project/internal/domain/models"
	"github.com/L0Qqi/go_final_project/internal/domain/services/nextDate"
)

func PostTaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func (app *App) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
		var task models.Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error": "Ошибка в декодировании JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error": "Поле title не может быть пустым"}`, http.StatusBadRequest)
			return
		}

		now := time.Now()
		if task.Date != "" {
			if _, err := time.Parse("20060102", task.Date); err != nil {
				http.Error(w, `{"error": "Неверный формат поля date"}`, http.StatusBadRequest)
				return
			}
		} else {
			task.Date = now.Format("20060102")
		}

		nextDate, err := nextDate.NextDateAdd(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка проверки даты: %v"}`, err), http.StatusBadRequest)
			return
		}
		task.Date = nextDate

		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"

		res, err := app.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка при добавлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка при получении id: %v"}`, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := fmt.Sprintf(`{"id": %d}`, id)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(response))
	}
}
