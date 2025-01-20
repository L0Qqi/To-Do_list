package taskhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/L0Qqi/To-Do_list/internal/app"
	"github.com/L0Qqi/To-Do_list/internal/domain/models"
)

// Возвращаем весь список задач отсортированных по дате
func GetTasksHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func (app *App) getTasksHandler(w http.ResponseWriter, r *http.Request) {

		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC"

		rows, err := app.DB.Query(query)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка выполнения запроса: %v"}`, err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tasks := []models.Task{}
		for rows.Next() {
			var id int
			var task models.Task
			if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "Ошибка обработки данных: %v"}`, err), http.StatusInternalServerError)
				return
			}
			task.ID = fmt.Sprintf("%d", id)
			tasks = append(tasks, task)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{
			"tasks": tasks,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Ошибка кодирования JSON: %v"}`, err), http.StatusInternalServerError)
		}
	}
}
