package routes

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type SqlLog struct {
	Id      int    `json:"id"`
	Cliente string `json:"cliente"`
	Numero  string `json:"numero"`
	Erro    string `json:"erro"`
	Data    string `json:"data"`
}

func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE log (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"cliente" TEXT,
		"numero" TEXT,
		"erro" TEXT,
		"data" DATETIME DEFAULT CURRENT_TIMESTAMP
	  );` // SQL Statement for Create Table

	statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("table created")
}

func InsertLog(cliente, numero string, erro string) {
	SQL := `INSERT INTO log(cliente, numero, erro) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(SQL)
	if err != nil {
		log.Println(err.Error())
	}
	_, err = statement.Exec(numero, erro)
	if err != nil {
		log.Println(err.Error())
	}
}
