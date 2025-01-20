package taskhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/L0Qqi/To-Do_list/internal/app"
	"github.com/L0Qqi/To-Do_list/internal/domain/models"
	"github.com/L0Qqi/To-Do_list/internal/domain/services"
)

// Обновляет задачу
func PutTaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func (app *App) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		//Декодируем тело запроса
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error": "Ошибка в декодировании JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error": "Поле title не может быть пустым"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			if _, err := time.Parse("20060102", task.Date); err != nil {
				http.Error(w, `{"error": "Неверный формат поля date"}`, http.StatusBadRequest)
				return
			}
		}
		//Проверяем правильность формата поля repeat
		if err := services.ValidateRepeat(task.Repeat); err != nil {
			http.Error(w, `{"error": "Неверный формат поля repeat"}`, http.StatusBadRequest)
			return
		}

		//Обновляем задачу
		query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"

		res, err := app.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Printf("Ошибка в проверке на : %v", err)
		}
		//Проверяет было ли затронуто 0 строк в результате запроса
		if rowsAffected == 0 {
			http.Error(w, `{"error": "Задача с указанным id не найдена или данные совпадают"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
