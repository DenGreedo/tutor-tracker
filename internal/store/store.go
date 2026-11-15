package store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotFound = errors.New("ученик не найден")

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Exam  string `json:"exam"`
	Price int    `json:"price"`
}

type Store struct {
	mu       sync.Mutex
	students []Student
	nextID   int
}

func (st *Store) FindByID(id int) (Student, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	for _, s := range st.students {
		if s.ID == id {
			return s, nil
		}
	}
	return Student{}, ErrNotFound
}

func (st *Store) Delete(id int) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	for i, s := range st.students {
		if s.ID == id {
			st.students = append(st.students[:i], st.students[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func NewStore() *Store {
	return &Store{nextID: 1}
}

func (st *Store) Add(name, exam string, price int) Student {
	st.mu.Lock()
	defer st.mu.Unlock()

	s := Student{ID: st.nextID, Name: name, Exam: exam, Price: price}
	st.students = append(st.students, s)
	st.nextID++
	return s
}

func (st *Store) Save(filename string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	data, err := json.MarshalIndent(st.students, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (st *Store) Load(filename string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	data, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &st.students)
	if err != nil {
		return err
	}

	maxID := 0
	for _, s := range st.students {
		if s.ID > maxID {
			maxID = s.ID
		}
	}
	st.nextID = maxID + 1

	return nil
}

func (st *Store) All() []Student {
	st.mu.Lock()
	defer st.mu.Unlock()

	result := make([]Student, len(st.students))
	copy(result, st.students)
	return result
}
