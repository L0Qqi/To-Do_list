package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/L0Qqi/go_final_project/database"
	"github.com/L0Qqi/go_final_project/nextDate"

	_ "modernc.org/sqlite"
)

type App struct {
	DB *sql.DB
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (app *App) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.postTaskHandler(w, r) // Обработка добавления задачи

	case http.MethodGet:
		app.getTaskHandler(w, r) // Обработка получения задачи

	case http.MethodPut: // Обработка изменения задачи
		app.putTaskHandler(w, r)

	// case http.MethodDelete:
	//     app.deleteTaskHandler(w, r) // Обработка удаления задачи

	default:
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}
}

func (app *App) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

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

	nextDate, err := nextDate.NextDate(now, task.Date, task.Repeat)
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

func (app *App) putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

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
	log.Printf("Rows affected: %d", rowsAffected)
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Задача с указанным id не найдена или данные совпадают"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (app *App) getTasksHandler(w http.ResponseWriter, r *http.Request) {

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC"

	rows, err := app.DB.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Ошибка выполнения запроса: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var id int
		var task Task
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

	jsonResponse, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonResponse)) // Временно для отладки

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Ошибка кодирования JSON: %v"}`, err), http.StatusInternalServerError)
	}
}

func (app *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1"

	row := app.DB.QueryRow(query, id)

	var task Task
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

func main() {

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	app := &App{DB: db}

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", nextDate.HandleNextDate)
	http.HandleFunc("/api/task", app.taskHandler)
	http.HandleFunc("/api/tasks", app.getTasksHandler)

	now := time.Now()
	fmt.Println(nextDate.NextDate(now, "20250105", "d 5"))

	if err := http.ListenAndServe(":7540", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
