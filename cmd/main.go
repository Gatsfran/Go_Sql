package main

import (
	"log"
	"net/http"
	"github.com/Gatsfran/Go_Sql/internal/config"
	"github.com/Gatsfran/Go_Sql/internal/controller"
	"github.com/Gatsfran/Go_Sql/internal/repo"
	
)

func main() {
	cfg := config.LoadConfig()

	db, err := repo.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	router := controller.NewRouter(db)
	

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}