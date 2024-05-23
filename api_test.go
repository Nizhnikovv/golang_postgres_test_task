package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

func (m *MockStorage) CreateEmployee(e Employee) error {
	return nil
}

func (m *MockStorage) GetEmployee(id int) (*Employee, error) {
	return nil, nil
}

func (m *MockStorage) GetEmployees() ([]*Employee, error) {
	return nil, nil
}

func TestCreateEmployee(t *testing.T) {
	rr := httptest.NewRecorder()

	employeeReq := CreateEmployeeRequest{
		FirstName:  "Dmitrii",
		LastName:   "Nizhnikov",
		MiddleName: "Victorovich",
		Phone:      "+77075143949",
		City:       "Almaty",
	}

	reqBody, err := json.Marshal(employeeReq)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Error(err)
	}

	mockStorage := &MockStorage{}
	server := APIServer{storage: mockStorage}
	server.createEmployee(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)

	var response CreateEmployeeRequest
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, employeeReq, response)
}
