package repo

import (
	"database/sql"
	"fmt"

	"github.com/Gatsfran/Go_Sql/internal/config"
	"github.com/Gatsfran/Go_Sql/internal/entity"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewDatabaseConnection(cfg config.Config) (*DB, error) {
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

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) GetReader(readerNum int) (*entity.Reader, error) {
	query := `
	SELECT 
		reader_num, 
		reader_name, 
		reader_adress, 
		reader_phone 
	FROM 
		readers 
	WHERE reader_num = $1`

	row := d.db.QueryRow(query, readerNum)
	reader := entity.Reader{}

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

func (d *DB) ListReader() ([]entity.Reader, error) {
	query := "SELECT reader_num, reader_name, reader_adress, reader_phone FROM readers"

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readers []entity.Reader
	for rows.Next() {
		var reader entity.Reader
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

func (d *DB) AddReader(reader entity.Reader) (int64, error) {
	sqlStatement := `
	INSERT INTO readers 
	(reader_name, reader_adress, reader_phone) 
	VALUES ($1, $2, $3) 
	RETURNING reader_num`

	var readerID int64
	err := d.db.QueryRow(sqlStatement, reader.Name, reader.Adress, reader.Phone, reader.ID).Scan(&readerID)

	return readerID, err
}
func(d *DB) UpdateReader (reader entity.Reader) error {
	query := `
	UPDATE readers SET 
		reader_name = $1, 
		reader_adress = $2, 
		reader_phone = $3 
	WHERE reader_num = $4`

	_, err := d.db.Exec(query, reader.Name, reader.Adress, reader.Phone, reader.ID)

	return err
}

func (d *DB) DeleteReader(readerNum int) error {
	query := `
	WITH
		(
			DELETE FROM books_in_use WHERE reader_num = $1
		)
	DELETE FROM readers WHERE reader_num = $1`

	_, err := d.db.Exec(query, readerNum)

	return err
}
