package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DenGreedo/tutor-tracker/internal/store"
)

func main() {
	const dataFile = "students.json"

	st := store.NewStore()
	if err := st.Load(dataFile); err != nil {
		log.Fatal("не удалось загрузить данные: ", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("GET /students", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(st.All()); err != nil {
			log.Println("не удалось отправить список учеников:", err)
		}
	})

	log.Println("сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
