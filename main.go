package main

import (
	"fmt"
	"net/http"
)

func main() {

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	if err := http.ListenAndServe(":7540", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
