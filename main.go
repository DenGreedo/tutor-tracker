package main

import (
	"log"
	"net/http"

	"github.com/DenGreedo/tutor-tracker/internal/httpapi"
	"github.com/DenGreedo/tutor-tracker/internal/store"
)

func main() {
	const dataFile = "students.json"

	st := store.NewStore()
	if err := st.Load(dataFile); err != nil {
		log.Fatal("не удалось загрузить данные: ", err)
	}

	log.Println("сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", httpapi.New(st, dataFile)))
}
