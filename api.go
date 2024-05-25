package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

type ApiError struct {
	Msg string `json:"error"`
}

func (e ApiError) Error() string {
	return e.Msg
}

func NewAPIServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()

	router.Get("/api/v1/employees", makeHTTPHandler(s.getEmployees))
	router.Get("/api/v1/employees/{id}", makeHTTPHandler(s.getEmployee))
	router.Post("/api/v1/employees", makeHTTPHandler(s.createEmployee))

	log.Printf("Starting server on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) getEmployees(w http.ResponseWriter, r *http.Request) error {
	employees, err := s.storage.GetEmployees()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, employees)
}

func (s *APIServer) getEmployee(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ApiError{Msg: fmt.Sprintf("invalid id given %s", idStr)}
	}

	employee, err := s.storage.GetEmployee(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, employee)
}

func (s *APIServer) createEmployee(w http.ResponseWriter, r *http.Request) error {
	var employeeReq CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&employeeReq); err != nil {
		return ApiError{Msg: "invalid request body"}
	}
	defer r.Body.Close()

	if err := s.storage.CreateEmployee(*employeeReq.ToEmployee()); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, employeeReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch err := fn(w, r); err.(type) {
		case nil:
			return
		case ApiError:
			WriteJSON(w, http.StatusBadRequest, ApiError{Msg: err.Error()})
		case StorageError:
			WriteJSON(w, http.StatusInternalServerError, ApiError{Msg: err.Error()})
		default:
			WriteJSON(w, http.StatusInternalServerError, ApiError{Msg: "internal server error"})
		}
	}
}
