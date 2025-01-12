package taskhandlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/L0Qqi/go_final_project/internal/app"
	"github.com/L0Qqi/go_final_project/internal/domain/models"
)

func validateRepeat(repeat string) error {
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

func PutTaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func (app *App) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
		var task models.Task

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

		if err := validateRepeat(task.Repeat); err != nil {
			http.Error(w, `{"error": "Неверный формат поля repeat"}`, http.StatusBadRequest)
			return
		}

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

		if rowsAffected == 0 {
			http.Error(w, `{"error": "Задача с указанным id не найдена или данные совпадают"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
