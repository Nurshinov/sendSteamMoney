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

func AddNewProduct(login string, amount string) string {
	var InsertId int
	err := MakeConn().QueryRow(`INSERT INTO transacts(login, amount)  
		VALUES ($1,$2) RETURNING id`, login, amount).Scan(&InsertId)
	if err != nil {
		log.Println("[ERROR] Неудачная запись в базу: " + err.Error())
	}
	return fmt.Sprintf("%.9b", InsertId)
}
