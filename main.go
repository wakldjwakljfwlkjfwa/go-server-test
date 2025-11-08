package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//go:embed index.html
var indexPage string

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	router := gin.Default()

	tmpl, err := template.New("index").Parse(indexPage)
	if err != nil {
		panic(err)
	}

	// Better use .env
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "mydb")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS employees (
        id SERIAL PRIMARY KEY,
        fullname TEXT,
        phone TEXT,
        city TEXT
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	var employeeRepo EmployeeRepositoryPostgres
	employeeRepo.Db = db

	// Test index page
	router.GET("/", func(ctx *gin.Context) {
		data := struct {
			Title string
			Body  string
		}{
			Title: "Hello, World!",
			Body:  "This is a dynamic HTML page served with Gin!",
		}

		if err := tmpl.Execute(ctx.Writer, data); err != nil {
			ctx.String(500, "Internal Server Error")
		}
	})

	router.GET("/employees", getEmployeesHandler(&employeeRepo))
	router.POST("/employees", postEmployeeHandler(&employeeRepo))

	router.Run(":8080")
}
