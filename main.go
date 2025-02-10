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

func ListReader(db *sql.DB) ([]Reader, error) {
	query := "SELECT reader_num, reader_name, reader_adress, reader_phone FROM readers"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readers []Reader
	for rows.Next() {
		var reader Reader
		err = rows.Scan(
			&reader.ID,
			&reader.Name,
			&reader.Adress,
			&reader.Phone,
		)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}

	return readers, nil
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

func AddReader(db *sql.DB, reader Reader) (int64, error) {
	sqlStatement := `
	INSERT INTO readers 
	(reader_name, reader_adress, reader_phone) 
	VALUES ($1, $2, $3) 
	RETURNING reader_num`

	var readerID int64
	err := db.QueryRow(sqlStatement, reader.Name, reader.Adress, reader.Phone, reader.ID).Scan(&readerID)
	return readerID, err
}
func UpdateReader(db *sql.DB, reader Reader) error {
	query := `
	UPDATE readers SET 
		reader_name = $1, 
		reader_adress = $2, 
		reader_phone = $3 
	WHERE reader_num = $4`
	_, err := db.Exec(query, reader.Name, reader.Adress, reader.Phone, reader.ID)
	return err
}

func DeleteReader(db *sql.DB, readerNum int) error {
	query := `DELETE FROM readers WHERE reader_num = $1`
	_, err := db.Exec(query, readerNum)
	return err
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

	addReader := Reader{
		Name:   "Вася",
		Adress: "Мира, 6",
		Phone:  "555545678",
	}
	readerID, err := AddReader(db, addReader)
	if err != nil {
		log.Printf("Ошибка при добавлении читателя: %v\n", err)
		return
	}
	fmt.Printf("Читатель успешно добавлен с ID: %d\n", readerID)

	updatedReader := Reader{
		ID:     1,
		Name:   "Вася",
		Adress: "Мира, 6",
		Phone:  "555545678",
	}
	err = UpdateReader(db, updatedReader)
	if err != nil {
		log.Printf("Ошибка при обновлении читателя: %v\n", err)
		return
	}

	reader, err := GetReader(db, 27)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(reader)

	readers, err := ListReader(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, reader := range readers {
		fmt.Printf("Номер: %d, Имя: %s, Адрес: %s, Телефон: %s\n",
			reader.ID,
			reader.Name,
			reader.Adress,
			reader.Phone,
		)
	}

	err = DeleteReader(db, 4)
	if err != nil {
		log.Printf("Ошибка при удалении читателя: %v\n", err)
		return
	}
}
