package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func MakeConn() *sql.DB {
	connSettings := fmt.Sprintf("host=%s dbname=%s sslmode=disable user=%s password=%s", os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	db, err := sql.Open("postgres",
		connSettings)
	if err != nil {
		log.Fatal("[ERROR] Cannot connect ...\n" + err.Error())
	}
	return db
}

func AddNewProduct(login string, amount string, operation string) string {
	var InsertId int
	err := MakeConn().QueryRow(`INSERT INTO transacts(login, amount,operation)  
		VALUES ($1,$2, $3) RETURNING id`, login, amount, operation).Scan(&InsertId)
	if err != nil {
		log.Println("[ERROR] Неудачная запись в базу: " + err.Error())
	}
	log.Printf("[INFO] Добавлена новая запись в базу id:%d, steamLogin: %s", InsertId, login)
	return fmt.Sprintf("%.14d", InsertId)
}
