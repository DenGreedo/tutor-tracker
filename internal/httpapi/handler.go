package httpapi

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

type handler struct {
	st       *store.Store
	dataFile string
}

func New(st *store.Store, dataFile string) *http.ServeMux {
	h := &handler{st: st, dataFile: dataFile}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /students", h.listStudents)
	mux.HandleFunc("GET /students/{id}", h.getStudent)
	mux.HandleFunc("POST /students", h.createStudent)
	mux.HandleFunc("DELETE /students/{id}", h.deleteStudent)
	return mux
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (h *handler) listStudents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.st.All())
}

func (h *handler) getStudent(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	s, err := h.st.FindByID(id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, s)
}

func (h *handler) createStudent(w http.ResponseWriter, r *http.Request) {
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

	s := h.st.Add(strings.TrimSpace(req.Name), req.Exam, req.Price)
	if !h.save(w) {
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

func (h *handler) deleteStudent(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if err := h.st.Delete(id); err != nil {
		writeStoreError(w, err)
		return
	}
	if !h.save(w) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) save(w http.ResponseWriter) bool {
	if err := h.st.Save(h.dataFile); err != nil {
		log.Println("не удалось сохранить данные:", err)
		writeError(w, http.StatusInternalServerError, "не удалось сохранить изменения")
		return false
	}
	return true
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id должен быть числом")
		return 0, false
	}
	return id, true
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "ученик не найден")
		return
	}
	log.Println("ошибка журнала:", err)
	writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
}

type createStudentRequest struct {
	Name  string `json:"name"`
	Exam  string `json:"exam"`
	Price int    `json:"price"`
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
