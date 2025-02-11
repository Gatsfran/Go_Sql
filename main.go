package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"

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

	r := mux.NewRouter()

	r.HandleFunc("/readers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			readers, err := ListReader(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(readers)
		case http.MethodPost:
			var reader Reader
			if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := AddReader(db, reader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Добавлен читатель с ID: %d", id)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	r.HandleFunc("/readers/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID читателя", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			reader, err := GetReader(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(reader)
		case http.MethodPut:
			var reader Reader
			if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reader.ID = id
			if err := UpdateReader(db, reader); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Читатель с ID %d обновлен", id)
		case http.MethodDelete:
			if err := DeleteReader(db, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Читатель с ID %d удален", id)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE")

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
	
}
