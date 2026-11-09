package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Exam  string `json:"exam"`
	Price int    `json:"price"`
}

type Store struct {
	students []Student
	nextID   int
}

func main() {

	const dataFile = "students.json"

	store := NewStore()
	if err := store.Load(dataFile); err != nil {
		fmt.Println("Не удалось загрузить данные:", err)
		return
	}
	store.Add("Карина", "ОГЭ", 1300)
	store.Add("Андрей", "ЕГЭ", 1500)
	store.Add("Настя", "ЕГЭ", 1500)
	store.Print()

	if err := store.Save(dataFile); err != nil {
		fmt.Println("Не удалось сохранить данные:", err)
	}

	s, err := store.FindByID(3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}

	s, err = store.FindByID(92)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}
}

func (st *Store) Print() {
	for _, s := range st.students {
		fmt.Printf("%d. %s - %s - %d ₽\n", s.ID, s.Name, s.Exam, s.Price)
	}
}

func (st *Store) FindByID(id int) (Student, error) {
	for _, s := range st.students {
		if s.ID == id {
			return s, nil
		}
	}
	return Student{}, errors.New("ученик не найден")
}

func NewStore() *Store {
	return &Store{nextID: 1}
}

func (st *Store) Add(name, exam string, price int) Student {
	s := Student{ID: st.nextID, Name: name, Exam: exam, Price: price}
	st.students = append(st.students, s)
	st.nextID++
	return s
}

func (st *Store) Save(filename string) error {
	data, err := json.MarshalIndent(st.students, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (st *Store) Load(filename string) error {
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
