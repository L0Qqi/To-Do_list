package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/L0Qqi/To-Do_list/database"
	"github.com/L0Qqi/To-Do_list/internal/app"
	"github.com/L0Qqi/To-Do_list/internal/taskhandlers"

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

	if err := http.ListenAndServe(":7540", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
