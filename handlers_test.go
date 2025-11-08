package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockEmployeeRepository struct{}

func (repo *MockEmployeeRepository) GetAllEmployees() ([]Employee, error) {
	var employees = []Employee{
		{
			Id:       1,
			Fullname: "John Doe",
			Phone:    "+7 777 707 77 77",
			City:     "Almaty",
		},
		{
			Id:       2,
			Fullname: "Alice",
			Phone:    "+7 707 123 45 67",
			City:     "Astana",
		},
		{
			Id:       3,
			Fullname: "Bob",
			Phone:    "+7 700 707 77 77",
			City:     "Almaty",
		},
	}
	return employees, nil
}

func (repo *MockEmployeeRepository) SaveEmployee(request *EmployeeRequest) (Employee, error) {
	var employee Employee
	employee.Fullname = request.Fullname
	employee.Phone = request.Phone
	employee.City = request.City
	employee.Id = 1
	return employee, nil
}

func TestGetEmployeesHandler(t *testing.T) {
	var mockRepo MockEmployeeRepository

	router := gin.Default()
	router.GET("/employees", getEmployeesHandler(&mockRepo))

	req := httptest.NewRequest("GET", "/employees", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Fatalf("Status codes do not match. Received %d    Expected %d", w.Code, http.StatusOK)
	}

	var response struct {
		Employees []Employee
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Employees) != 3 {
		t.Fatalf("Expected 3 employees, but got %d", len(response.Employees))
	}

	if response.Employees[0].Id != 1 {
		t.Fatalf("Expected employee ID 1, but got %d", response.Employees[0].Id)
	}
}

func TestPostEmployeeHandler(t *testing.T) {
	var mockRepo MockEmployeeRepository

	router := gin.Default()
	router.POST("/employees", postEmployeeHandler(&mockRepo))

	request := EmployeeRequest{Fullname: "John Doe", Phone: "+7 777 777 77 77", City: "Almaty"}
	requestJson, _ := json.Marshal(request)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/employees", bytes.NewBuffer(requestJson))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if http.StatusCreated != w.Code {
		t.Fatalf("Status codes do not match. Received %d    Expected %d", w.Code, http.StatusCreated)
	}

	var employee Employee

	err := json.Unmarshal(w.Body.Bytes(), &employee)
	if err != nil {
		t.Fatalf("Failed to unmarshal response into employee: %v", err)
	}

	if employee.Id != 1 {
		t.Fatalf("Expected employee ID to be 1, but got %d", employee.Id)
	}
}
