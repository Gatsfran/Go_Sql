package main

import (
	"log"
	"net/http"

	"github.com/Gatsfran/Go_Sql/internal/config"
	"github.com/Gatsfran/Go_Sql/internal/controller"
	"github.com/Gatsfran/Go_Sql/internal/repo"
)

func main() {
	cfg := config.NewCfg()

	db, err := repo.New(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	router := controller.New(db)

	serverAddr := ":" + cfg.Server.Port
	log.Printf("Сервер запущен на %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
