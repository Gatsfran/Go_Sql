package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5400
	user     = "postgres"
	password = "docker"
	dbname   = "postgres"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	defer db.Close()

	insertStmt := `insert into "readers" ("reader_name", "reader_adress", "reader_phone") 
    values('Гаценко', 'Зеленая, 3', 1111000)`
	_, e := db.Exec(insertStmt)
	CheckError(e)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
