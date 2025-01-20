package taskhandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/L0Qqi/To-Do_list/internal/app"
	"github.com/L0Qqi/To-Do_list/internal/domain/models"
)

// Ищем задачу по её id
func GetTaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//func (app *App) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1"

		row := app.DB.QueryRow(query, id)

		var task models.Task
		var taskID int

		if err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"Задача с указаным id не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка выполнения запроса: %v"}`, err), http.StatusInternalServerError)
			}
			return
		}
		task.ID = fmt.Sprintf("%d", taskID)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка кодирования JSON: %v"}`, err), http.StatusInternalServerError)
		}
	}
}
