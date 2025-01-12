package taskhandlers

import (
	"net/http"

	"github.com/L0Qqi/go_final_project/internal/app"
)

// func (app *App) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

func DeleteTaskHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"ID задачи не указан"}`, http.StatusBadRequest)
			return
		}

		query := "DELETE FROM scheduler WHERE id = $1"

		res, err := app.DB.Exec(query, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
