package entity

import "fmt"

type Reader struct {
	ID     int     `json:"num"`
	Name   string  `json:"name"`
	Adress *string `json:"adress"`
	Phone  string  `json:"phone"`
}

func (r Reader) String() string {
	return fmt.Sprintf(`
	Информация о читателе:
	========================
	Номер читателя: %d
	Имя читателя:  %s
	Адрес:        %s
	Телефон:      %s
	========================
	`, r.ID, r.Name, *r.Adress, r.Phone)
}