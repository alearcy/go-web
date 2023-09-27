package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// Db - the db instance
var Db *sql.DB

const statment string = `
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER,
    uuid TEXT NOT NULL UNIQUE, 
    user_id INTEGER NOT NULL UNIQUE, 
    created_at DATETIME NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role INT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);`

// StartDb - db initialization function
func StartDb() {
	viper.SetConfigFile(".env")
	var err error
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// dbName, _ := viper.Get("DB_NAME").(string)
	// dbUser, _ := viper.Get("DB_USER").(string)
	// dbHost, _ := viper.Get("DB_HOST").(string)
	// dbPwd, _ := viper.Get("DB_PWD").(string)
	// dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPwd, dbHost, dbName)
	Db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	_, err = Db.Exec(statment)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the DB")
}
