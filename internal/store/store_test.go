package store

import (
	"path/filepath"
	"testing"
)

func TestAdd_AssignsSequentialIDs(t *testing.T) {
	store := NewStore()

	first := store.Add("Карина", "ОГЭ", 1300)
	second := store.Add("Андрей", "ЕГЭ", 1500)

	if first.ID != 1 {
		t.Errorf("первый ученик: ID = %d, ожидали 1", first.ID)
	}

	if second.ID != 2 {
		t.Errorf("второй ученик: ID = %d, ожидали 2", second.ID)
	}
}

func TestFindByID_Found(t *testing.T) {
	store := NewStore()

	store.Add("Карина", "ОГЭ", 1500)
	store.Add("Андрей", "ЕГЭ", 1500)
	store.Add("Настя", "ОГЭ", 1500)

	s, err := store.FindByID(2)

	if err != nil {
		t.Fatalf("FindByID(2) вернул ошибку: %v", err)
	}

	if s.Name != "Андрей" {
		t.Errorf("FindByID(2): Name = %s, ожидали Андрей", s.Name)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	store := NewStore()

	store.Add("Карина", "ОГЭ", 1500)
	store.Add("Андрей", "ЕГЭ", 1500)
	store.Add("Настя", "ОГЭ", 1500)

	s, err := store.FindByID(99)

	if err == nil {
		t.Fatal("FindByID(99) не вернул ошибку, ученика с ID 99 нет")
	}

	if s != (Student{}) {
		t.Errorf("FindByID(99) вернул ученика %+v, ожидали пустоту", s)
	}

}

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "students.json")

	original := NewStore()
	original.Add("Карина", "ОГЭ", 1300)
	original.Add("Андрей", "ЕГЭ", 1500)

	if err := original.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded := NewStore()

	if err := loaded.Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.students) != 2 {
		t.Errorf("загружено учеников: %d, ожидали 2", len(loaded.students))
	}

	if loaded.nextID != 3 {
		t.Errorf("nextID после загрузки = %d, ожидали 3", loaded.nextID)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	store := NewStore()
	store.Add("Карина", "ОГЭ", 1300)
	store.Add("Андрей", "ЕГЭ", 1500)

	result := store.All()

	if len(result) != 2 {
		t.Fatalf("All вернул %d учеников, ожидали 2", len(result))
	}

	result[0].Name = "Изменено"

	s, err := store.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID(1) вернул ошибку: %v", err)
	}
	if s.Name != "Карина" {
		t.Errorf("имя в журнале = %s, ожидали Карина: All отдал не копию", s.Name)
	}
}
