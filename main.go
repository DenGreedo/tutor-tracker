package main

import (
	"errors"
	"fmt"
)

type Student struct {
	ID    int
	Name  string
	Exam  string
	Price int
}

type Store struct {
	students []Student
	nextID   int
}

func main() {

	store := NewStore()
	store.Add("Карина", "ОГЭ", 1300)
	store.Add("Андрей", "ЕГЭ", 1500)
	store.Add("Настя", "ЕГЭ", 1500)
	store.Print()

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
