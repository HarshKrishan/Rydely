package db

import (
	"context"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var client *mongo.Client
var collection *mongo.Collection

func InitMongoDB() {

	host := os.Getenv("MONGO_URI")

	rb := bson.NewRegistryBuilder()

	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))

	clientOptions := options.Client().ApplyURI(host).SetRegistry(rb.Build())

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Panicln("Error connecting to MongoDB:", err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Panicln("Error pinging MongoDB:", err)
	}

	db = client.Database("rydely")
	log.Infoln("Connected to Mongo DB!")
}

func GetDB() *mongo.Database {
	if db == nil {
		log.Panicln("MongoDB is not initialized. Call InitMongoDB() first.")
	}

	return db
}

func GetCollection() *mongo.Collection {
	if db == nil {
		log.Panicln("MongoDB is not initialized. Call InitMongoDB() first.")
	}

	return GetDB().Collection("rides")
}
