package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Gatsfran/Go_Sql/internal/entity"
	"github.com/Gatsfran/Go_Sql/internal/repo"
	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
	db     *repo.DB
}

func New(db *repo.DB) *Router {
	r := mux.NewRouter()
	router := &Router{
		router: r,
		db:     db,
	}
	router.newRegisterReaderRoutes()
	return router
}

func (r *Router) newRegisterReaderRoutes() {
	r.router.HandleFunc("/readers", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			readers, err := r.db.ListReader()
			if err != nil {
				log.Printf("Ошибка при получении списка читателей: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.NewEncoder(w).Encode(readers); err != nil {
				log.Printf("Ошибка при кодировании списка читателей в JSON: %v", err)
				return
			}

		case http.MethodPost:
			var reader entity.Reader
			if err := json.NewDecoder(req.Body).Decode(&reader); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := r.db.AddReader(reader)
			if err != nil {
				log.Printf("Ошибка при добавлении читателя: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			if _, err := fmt.Fprintf(w, "Добавлен читатель с ID: %d", id); err != nil {
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	r.router.HandleFunc("/readers/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Printf("Неверный ID читателя: %v", err)
			http.Error(w, "Неверный ID читателя", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodGet:
			reader, err := r.db.GetReader(id)
			if err != nil {
				log.Printf("Ошибка при получении читателя с ID %d: %v", id, err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			if err := json.NewEncoder(w).Encode(reader); err != nil {
				log.Printf("Ошибка при кодировании читателя в JSON: %v", err)
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		case http.MethodPut:
			var reader entity.Reader
			if err := json.NewDecoder(req.Body).Decode(&reader); err != nil {
				log.Printf("Ошибка при декодировании тела запроса: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reader.ID = id
			if err := r.db.UpdateReader(reader); err != nil {
				log.Printf("Ошибка при обновлении читателя с ID %d: %v", id, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := fmt.Fprintf(w, "Читатель с ID %d обновлен", id); err != nil {
				log.Printf("Ошибка при записи ответа: %v", err)
				http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
				return
			}

		case http.MethodDelete:
			if err := r.db.DeleteReader(id); err != nil {
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
