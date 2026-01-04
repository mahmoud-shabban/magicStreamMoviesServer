package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: can not find .env file")
	}

	mongodb, ok := os.LookupEnv("MONGODB_URI")

	if !ok {
		log.Fatal("MONGODB_URI not set")
	}

	clientOptions := options.Client().ApplyURI(mongodb)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Println("Error: cannot connect to db ", err.Error())
		return nil
	}

	return client
}

var Client *mongo.Client = Connect()

func OpenCollection(collectionName string) *mongo.Collection {

	dbName, ok := os.LookupEnv("DATABASE_NAME")

	if !ok {
		log.Println("DATABASE_NAME not set")
		return nil
	}

	if Client == nil {
		log.Fatal("cannot connect to db")
	}

	return Client.Database(dbName).Collection(collectionName)
}
