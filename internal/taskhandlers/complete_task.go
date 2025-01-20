package taskhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/L0Qqi/To-Do_list/internal/app"
	"github.com/L0Qqi/To-Do_list/internal/domain/models"
	"github.com/L0Qqi/To-Do_list/internal/domain/services/nextDate"
)

// Выполнение задачи
func TaskDoneHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error":"Не указан id задачи"}`, http.StatusBadRequest)
			return
		}

		//Считываем дату и правило повторения
		var task models.Task
		query := "SELECT date, repeat FROM scheduler WHERE id = $1"

		err := app.DB.QueryRow(query, id).Scan(&task.Date, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error":"Ошибка чтения задачи из базы данных"}`, http.StatusInternalServerError)
			}
			return
		}

		// Если повторения нет, удаляем задачу
		if task.Repeat == "" {
			query := "DELETE FROM scheduler WHERE id = $1"
			_, err := app.DB.Exec(query, id)
			if err != nil {
				http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{}`))
			return
		}

		// Рассчитываем следующую дату
		now := time.Now()
		nextDate, err := nextDate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка расчета следующей даты: %v"}`, err), http.StatusBadRequest)
			return
		}

		// Обновляем дату в базе
		query = "UPDATE scheduler SET date = $1 WHERE id = $2"
		_, err = app.DB.Exec(query, nextDate, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}
}
