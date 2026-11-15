package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DenGreedo/tutor-tracker/internal/store"
)

func newTestMux(t *testing.T) *http.ServeMux {
	t.Helper()

	st := store.NewStore()
	st.Add("Карина", "ОГЭ", 1300)
	st.Add("Андрей", "ЕГЭ", 1500)

	path := filepath.Join(t.TempDir(), "students.json")
	return New(st, path)
}

func doRequest(mux http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	mux := newTestMux(t)

	rec := doRequest(mux, http.MethodGet, "/health", "")

	if rec.Code != http.StatusOK {
		t.Errorf("код = %d, ожидали 200", rec.Code)
	}
	if body := rec.Body.String(); body != "ok\n" {
		t.Errorf("тело = %q, ожидали %q", body, "ok\n")
	}
}

func TestGetStudents(t *testing.T) {
	mux := newTestMux(t)

	rec := doRequest(mux, http.MethodGet, "/students", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("код = %d, ожидали 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, ожидали application/json", ct)
	}

	var got []store.Student
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("ответ — не JSON: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("учеников в ответе = %d, ожидали 2", len(got))
	}
}

func TestGetStudentByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantCode int
	}{
		{"существующий ученик", "1", http.StatusOK},
		{"несуществующий ученик", "999", http.StatusNotFound},
		{"id не число", "abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := newTestMux(t)

			rec := doRequest(mux, http.MethodGet, "/students/"+tt.id, "")

			if rec.Code != tt.wantCode {
				t.Errorf("код = %d, ожидали %d; тело: %s", rec.Code, tt.wantCode, rec.Body.String())
			}
		})
	}
}

func TestCreateStudent(t *testing.T) {
	mux := newTestMux(t)

	rec := doRequest(mux, http.MethodPost, "/students",
		`{"name":"Петя","exam":"ЕГЭ","price":1600}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("код = %d, ожидали 201; тело: %s", rec.Code, rec.Body.String())
	}

	var got store.Student
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("ответ — не JSON: %v", err)
	}
	if got.ID != 3 || got.Name != "Петя" {
		t.Errorf("создан ученик %+v, ожидали Петю с ID 3", got)
	}

	rec = doRequest(mux, http.MethodGet, "/students/3", "")
	if rec.Code != http.StatusOK {
		t.Errorf("GET /students/3 после создания: код = %d, ожидали 200", rec.Code)
	}
}

func TestCreateStudent_Invalid(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"пустое имя", `{"name":"","exam":"ЕГЭ","price":1500}`},
		{"имя из пробелов", `{"name":"   ","exam":"ЕГЭ","price":1500}`},
		{"неизвестный экзамен", `{"name":"Петя","exam":"ВПР","price":1500}`},
		{"отрицательная цена", `{"name":"Петя","exam":"ЕГЭ","price":-100}`},
		{"сломанный JSON", `{"name":"Петя"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := newTestMux(t)

			rec := doRequest(mux, http.MethodPost, "/students", tt.body)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("код = %d, ожидали 400; тело: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestDeleteStudent(t *testing.T) {
	mux := newTestMux(t)

	rec := doRequest(mux, http.MethodDelete, "/students/1", "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("первое удаление: код = %d, ожидали 204", rec.Code)
	}

	rec = doRequest(mux, http.MethodDelete, "/students/1", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("повторное удаление: код = %d, ожидали 404", rec.Code)
	}

	rec = doRequest(mux, http.MethodGet, "/students/1", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("GET после удаления: код = %d, ожидали 404", rec.Code)
	}
}
