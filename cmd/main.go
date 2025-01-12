package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/L0Qqi/go_final_project/database"
	"github.com/L0Qqi/go_final_project/internal/app"
	"github.com/L0Qqi/go_final_project/internal/domain/services/nextDate"
	"github.com/L0Qqi/go_final_project/internal/taskhandlers"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	app := &app.App{DB: db}

	webDir := "../web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", taskhandlers.HandleNextDate)
	http.HandleFunc("/api/task", taskhandlers.TaskHandler(app))
	http.HandleFunc("/api/tasks", taskhandlers.GetTasksHandler(app))
	http.HandleFunc("/api/task/done", taskhandlers.TaskDoneHandler(app))

	now := time.Now()
	fmt.Println(nextDate.NextDate(now, "20250111", "d 5"))

	if err := http.ListenAndServe(":7540", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
