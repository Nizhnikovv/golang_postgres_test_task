package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetEmployees() ([]*Employee, error)
	GetEmployee(id int) (*Employee, error)
	CreateEmployee(employee Employee) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres password=12345678 dbname=gopost sslmode=disable host=172.17.0.5"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	query := `
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			middle_name TEXT NOT NULL,
			phone TEXT NOT NULL,
			city TEXT NOT NULL
		);
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateEmployee(employee Employee) error {
	query := `
		INSERT INTO employees (first_name, last_name, middle_name, phone, city)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(query, employee.FirstName, employee.LastName, employee.MiddleName, employee.Phone, employee.City)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) GetEmployees() ([]*Employee, error) {
	rows, err := s.db.Query("SELECT id, first_name, last_name, middle_name, phone, city FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []*Employee{}
	for rows.Next() {
		employee, err := scanIntoEmployee(rows)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func (s *PostgresStorage) GetEmployee(id int) (*Employee, error) {
	rows, err := s.db.Query("SELECT * FROM employees WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoEmployee(rows)
	}

	return nil, fmt.Errorf("employee with id %d not found", id)
}

func scanIntoEmployee(rows *sql.Rows) (*Employee, error) {
	employee := new(Employee)
	if err := rows.Scan(
		&employee.ID,
		&employee.FirstName,
		&employee.LastName,
		&employee.MiddleName,
		&employee.Phone,
		&employee.City,
	); err != nil {
		return nil, err
	}

	return employee, nil
}
