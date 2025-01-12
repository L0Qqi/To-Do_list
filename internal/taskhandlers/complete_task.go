package taskhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/L0Qqi/go_final_project/internal/app"
	"github.com/L0Qqi/go_final_project/internal/domain/models"
	"github.com/L0Qqi/go_final_project/internal/domain/services/nextDate"
)

func TaskDoneHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//func (app *app.App) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

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

		flag := (task.Repeat == "")

		switch flag {
		case true:
			query := "DELETE FROM scheduler WHERE id = $1"

			_, err := app.DB.Exec(query, id)
			if err != nil {
				http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{}`))

		case false:
			now := time.Now()
			nextDate, err := nextDate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Ошибка расчета следующей даты"}`, http.StatusBadRequest)
				return
			}
			query := "UPDATE scheduler SET date = $1 WHERE id = $2"

			_, err = app.DB.Exec(query, nextDate, id)
			if err != nil {
				fmt.Printf("Ошибка обновления задачи: %v\n", err)
				http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
			fmt.Printf("Обновлённая дата сохранена в базе: %s\n", nextDate)

			w.Header().Set("Content-Type", "application/json")

			w.Write([]byte(`{}`))
		}
	}
}
