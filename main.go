package main

import (
	"fmt"

	"github.com/DenGreedo/tutor-tracker/internal/store"
)

func main() {

	const dataFile = "students.json"

	st := store.NewStore()
	if err := st.Load(dataFile); err != nil {
		fmt.Println("Не удалось загрузить данные:", err)
		return
	}
	st.Add("Карина", "ОГЭ", 1300)
	st.Add("Андрей", "ЕГЭ", 1500)
	st.Add("Настя", "ЕГЭ", 1500)
	st.Print()

	s, err := st.FindByID(3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}

	s, err = st.FindByID(92)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Найден: %+v\n", s)
	}

	if err := st.Save(dataFile); err != nil {
		fmt.Println("Не удалось сохранить данные:", err)
	}
}
