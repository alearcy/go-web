package database

import (
	"context"
	_ "database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"

	"github.com/spf13/viper"
)

// CallDb - db initialization function
func CallDb() *mongo.Database {
	viper.SetConfigFile(".env")
	//var err error
	//var cancel context.CancelFunc
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	dbHost, _ := viper.Get("DB_HOST").(string)
	dbName, _ := viper.Get("DB_NAME").(string)
	client, err := mongo.NewClient(options.Client().ApplyURI(dbHost))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(Ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(dbName)
}
