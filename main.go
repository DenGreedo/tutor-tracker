package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DenGreedo/tutor-tracker/internal/store"
)

type createStudentRequest struct {
	Name  string `json:"name"`
	Exam  string `json:"exam"`
	Price int    `json:"price"`
}

func main() {
	const dataFile = "students.json"

	st := store.NewStore()
	if err := st.Load(dataFile); err != nil {
		log.Fatal("не удалось загрузить данные: ", err)
	}

	mux := newMux(st, dataFile)

	log.Println("сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func newMux(st *store.Store, dataFile string) *http.ServeMux {
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

	mux.HandleFunc("POST /students", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var req createStudentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "некорректный JSON")
			return
		}

		if err := validateStudent(req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		s := st.Add(strings.TrimSpace(req.Name), req.Exam, req.Price)

		if err := st.Save(dataFile); err != nil {
			log.Println("не удалось сохранить данные:", err)
			writeError(w, http.StatusInternalServerError, "не удалось сохранить ученика")
			return
		}

		writeJSON(w, http.StatusCreated, s)
	})

	mux.HandleFunc("DELETE /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "id должен быть числом")
			return
		}

		err = st.Delete(id)
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
			return
		}

		if err := st.Save(dataFile); err != nil {
			log.Println("не удалось сохранить данные:", err)
			writeError(w, http.StatusInternalServerError, "не удалось сохранить изменения")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	return mux
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

func validateStudent(req createStudentRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("имя не может быть пустым")
	}
	if req.Exam != "ЕГЭ" && req.Exam != "ОГЭ" {
		return errors.New("exam должен быть ЕГЭ или ОГЭ")
	}
	if req.Price <= 0 {
		return errors.New("цена должна быть больше нуля")
	}
	return nil
}
