package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Db - the db instance
var Db *sql.DB

// StartDb - db initialization function
func StartDb() {
	viper.SetConfigFile(".env")
	var err error
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	dbName, _ := viper.Get("DB_NAME").(string)
	dbUser, _ := viper.Get("DB_USER").(string)
	dbHost, _ := viper.Get("DB_HOST").(string)
	dbPwd, _ := viper.Get("DB_PWD").(string)
	dbPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPwd, dbHost, dbName)
	Db, err = sql.Open("pgx", dbPath)
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the DB")
}
