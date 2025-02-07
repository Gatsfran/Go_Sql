package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Reader struct {
	ID     int    `json:"num"`
	Name   string `json:"name"`
	Adress string `json:"adress"`
	Phone  string `json:"phone"`
}

func NewDatabaseConnection(cfg Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии соединения с БД: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к БД: %w", err)
	}

	return db, nil
}

func GetReader(db *sql.DB, readerNum int) (*Reader, error) {
	query := `
	SELECT 
		reader_num, 
		reader_name, 
		reader_adress, 
		reader_phone 
	FROM 
		readers 
	WHERE reader_num = $1`

	row := db.QueryRow(query, readerNum)
	reader := Reader{}

	err := row.Scan(
		&reader.ID,
		&reader.Name,
		&reader.Adress,
		&reader.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("читатель с номером %d не найден", readerNum)
		}
		return nil, fmt.Errorf("ошибка при чтении данных читателя: %w", err)
	}

	return &reader, nil
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
	`, r.ID, r.Name, r.Adress, r.Phone)
	}

func main() {
	config := Config{
		Host:     "localhost",
		Port:     "5400",
		Username: "postgres",
		Password: "docker",
		Database: "postgres",
	}

	db, err := NewDatabaseConnection(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	reader, err := GetReader(db, 2)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(reader)

}
