package main

import (
    "database/sql"
    "fmt"

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
    ReaderNum   int    `json:"reader_num"`
    ReaderName  string `json:"reader_name"`
    ReaderAdress string `json:"reader_adress"`
    ReaderPhone string `json:"reader_phone"`
}


func NewDatabaseConnection(cfg Config) (*sql.DB, error) {
    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
    
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("ошибка при открытии соединения с БД: %v", err)
    }
    
    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("ошибка при проверке подключения к БД: %v", err)
    }
    
    return db, nil
}

func GetReader(db *sql.DB, readerNum int) (*Reader, error) {
    query := "SELECT reader_num, reader_name, reader_adress, reader_phone FROM readers WHERE reader_num = $1"
    
    row := db.QueryRow(query, readerNum)
    reader := &Reader{}
    
    err := row.Scan(&reader.ReaderNum, &reader.ReaderName, &reader.ReaderAdress, &reader.ReaderPhone)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("читатель с номером %d не найден", readerNum)
        }
        return nil, fmt.Errorf("ошибка при чтении данных читателя: %v", err)
    }
    
    return reader, nil
}


func PrintReader(reader *Reader) {
    fmt.Printf("\nИнформация о читателе:\n")
    fmt.Printf("========================\n")
    fmt.Printf("Номер читателя: %d\n", reader.ReaderNum)
    fmt.Printf("Имя читателя:  %s\n", reader.ReaderName)
    fmt.Printf("Адрес:        %s\n", reader.ReaderAdress)
    fmt.Printf("Телефон:      %s\n", reader.ReaderPhone)
    fmt.Printf("========================\n")
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
        fmt.Println(err)
        return
    }
    defer db.Close()

    
    reader, err := GetReader(db, 1)
    if err != nil {
        fmt.Println(err)
        return
    }
    PrintReader(reader)


    insertStmt := `insert into "readers" ("reader_name", "reader_adress", "reader_phone") 
    values('Гаценко', 'Зеленая, 3', 1111000)`
    fmt.Println("Добавлен новый читатель: ", insertStmt)

   

}

