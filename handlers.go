package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	Id       int    `json:"id"`
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
}

type EmployeeRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	City     string `json:"city" binding:"required"`
}

type EmployeeRepository interface {
	GetAllEmployees() ([]Employee, error)
	SaveEmployee(request *EmployeeRequest) (Employee, error)
}

type EmployeeRepositoryPostgres struct {
	Db *sql.DB
}

func (repo *EmployeeRepositoryPostgres) GetAllEmployees() ([]Employee, error) {
	employees := []Employee{}

	rows, err := repo.Db.Query("SELECT id, fullname, phone, city FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.Id, &employee.Fullname, &employee.Phone, &employee.City); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (repo *EmployeeRepositoryPostgres) SaveEmployee(request *EmployeeRequest) (Employee, error) {
	var employee Employee
	var lastInsertId int
	insertQuery := `INSERT INTO employees (fullname, phone, city) VALUES ($1, $2, $3) RETURNING id`
	err := repo.Db.QueryRow(insertQuery, request.Fullname, request.Phone, request.City).Scan(&lastInsertId)
	if err != nil {
		return employee, err
	}

	employee.Phone = request.Phone
	employee.Fullname = request.Fullname
	employee.City = request.City
	employee.Id = lastInsertId
	return employee, nil
}

func getEmployeesHandler(repo EmployeeRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		employees, err := repo.GetAllEmployees()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Show errors for testing
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"employees": employees,
		})
	}
}

func postEmployeeHandler(repo EmployeeRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request EmployeeRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Show errors for testing
			return
		}

		employee, err := repo.SaveEmployee(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Show errors for testing
			return
		}

		ctx.JSON(http.StatusCreated, employee)
	}
}
