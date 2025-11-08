# Test Go Web Server

## Task

Техническое задание для backend (GoLang):
Необходимо написать веб-сервис, который записывает данные сотрудника (ФИО, телефон, город) в базу PostgreSQL.
Всё оберните в Докер. Необходимо написать парочку тестов.

## How to run

```sh
docker compose up
```

The server listens to port 8080

API contains following routes:
```
GET http://localhost:8080/employees
POST http://localhost:8080/employees
```

For testing purposes I created form page
```
GET http://localhost:8080/
```
