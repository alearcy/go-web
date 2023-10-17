package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Db - the db instance
var Db *sql.DB

const statment string = `
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE, 
    user_id INTEGER NOT NULL UNIQUE, 
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role INT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
	
CREATE TABLE IF NOT EXISTS brothers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
    surname TEXT NOT NULL,
	gender TEXT NOT NULL,
	is_married BOOLEAN NOT NULL,
	spouse_id INTEGER REFERENCES brothers (id),
	not_available_days TEXT,
	is_pioneer BOOLEAN NOT NULL,
	is_sm BOOLEAN NOT NULL,
	is_az BOOLEAN NOT NULL,
	is_available BOOLEAN NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS program (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	month TEXT NOT NULL,
	day TEXT NOT NULL,
	time TEXT NOT NULL,
	address TEXT NOT NULL,
	b1_id INTEGER NOT NULL, 
	b2_id INTEGER NOT NULL, 
	b3_id INTEGER, 
	b4_id INTEGER, 
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,	
	FOREIGN KEY (b1_id, b2_id, b3_id, b4_id) REFERENCES brothers
);`

// StartDb - db initialization function
func StartDb() {
	//viper.SetConfigFile("../../.env")
	var err error
	//err = viper.ReadInConfig()
	//if err != nil {
	//	log.Fatalf("Error while reading config file %s", err)
	//}
	// dbName, _ := viper.Get("DB_NAME").(string)
	// dbUser, _ := viper.Get("DB_USER").(string)
	// dbHost, _ := viper.Get("DB_HOST").(string)
	// dbPwd, _ := viper.Get("DB_PWD").(string)
	// dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPwd, dbHost, dbName)
	Db, err = sql.Open("sqlite3", "../../database.db")
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
