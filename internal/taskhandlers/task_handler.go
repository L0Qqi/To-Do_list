package taskhandlers

import (
	"net/http"

	"github.com/L0Qqi/go_final_project/internal/app"
)

func TaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func (app *App) taskHandler(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PostTaskHandler(app)(w, r) // Обработка добавления задачи

		case http.MethodGet:
			GetTaskHandler(app)(w, r) // Обработка получения задачи

		case http.MethodPut: // Обработка изменения задачи
			PutTaskHandler(app)(w, r)

		case http.MethodDelete:
			DeleteTaskHandler(app)(w, r) // Обработка удаления задачи

		default:
			http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	}
}
