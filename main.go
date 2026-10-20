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

func main() {
	students := []Student{
		{ID: 1, Name: "Карина", Exam: "ОГЭ", Price: 1300},
		{ID: 2, Name: "Андрей", Exam: "ЕГЭ", Price: 1500},
		{ID: 3, Name: "Настя", Exam: "ЕГЭ", Price: 1500},
	}

	printStudents(students)

	s, err := findByID(students, 2)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}

	s, err = findByID(students, 92)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}
}

func printStudents(students []Student) {
	for _, s := range students {
		fmt.Printf("%d. %s - %s - %d ₽\n", s.ID, s.Name, s.Exam, s.Price)
	}
}

func findByID(students []Student, id int) (Student, error) {
	for _, s := range students {
		if s.ID == id {
			return s, nil
		}
	}
	return Student{}, errors.New("ученик не найден")
}
