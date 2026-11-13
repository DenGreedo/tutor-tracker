package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		writeJSON(w, http.StatusOK, st.All())
	})

	mux.HandleFunc("GET /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "id должен быть числом")
			return
		}

		s, err := st.FindByID(id)
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}

		writeJSON(w, http.StatusOK, s)
	})

	log.Println("сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("не удалось записать JSON:", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
