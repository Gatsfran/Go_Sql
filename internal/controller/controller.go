package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/Gatsfran/Go_Sql/internal/entity"
	"github.com/Gatsfran/Go_Sql/internal/repo"
	"github.com/gorilla/mux"
)

func RegisterReaderRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/readers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			readers, err := repo.ListReader(db)
			if err != nil {
				log.Printf("Ошибка при получении списка читателей: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(readers)

		case http.MethodPost:
			var reader entity.Reader
			if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := repo.AddReader(db, reader)
			if err != nil {
				log.Printf("Ошибка при добавлении читателя: %v", err)
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
			log.Printf("Неверный ID читателя: %v", err)
			http.Error(w, "Неверный ID читателя", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			reader, err := repo.GetReader(db, id)
			if err != nil {
				log.Printf("Ошибка при получении читателя с ID %d: %v", id, err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(reader)

		case http.MethodPut:
			var reader entity.Reader
			if err := json.NewDecoder(r.Body).Decode(&reader); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reader.ID = id
			if err := repo.UpdateReader(db, reader); err != nil {
				log.Printf("Ошибка при обновлении читателя с ID %d: %v", id, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Читатель с ID %d обновлен", id)

		case http.MethodDelete:
			if err := repo.DeleteReader(db, id); err != nil {
				log.Printf("Ошибка при удалении читателя с ID %d: %v", id, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Читатель с ID %d удален", id)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "PUT", "DELETE")
}
