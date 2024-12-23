package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/L0Qqi/go_final_project/database"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer db.Close()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":7540", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
