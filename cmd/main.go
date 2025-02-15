package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gatsfran/Go_Sql/internal/config"
	"github.com/Gatsfran/Go_Sql/internal/controller"
	"github.com/Gatsfran/Go_Sql/internal/repo"
)

func main() {
	cfg := config.New()

	db, err := repo.New(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	router := controller.New(db)

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Сервер запущен на порту: %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
