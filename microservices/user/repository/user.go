package repository

import (
	"context"
	"user/db"
	"user/models"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct{}

func (ur UserRepository) CheckUserExists(email string) (bool, error) {

	collection := db.GetCollection()

	var filter interface{}
	response := models.User{}
	filter = bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// User does not exist
			return false, nil
		}
		log.Errorln("Error finding user:", err)
		return false, err
	}

	// Decode the user document into response
	err := result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding user:", err)
		return false, err
	}

	return true, nil
}

func (ur UserRepository) CreateUser(user models.User) error {
	collection := db.GetCollection()

	// Insert the user document into the collection
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Errorln("Error inserting user:", err)
		return err
	}

	log.Infoln("User created successfully")
	return nil
}

func (ur UserRepository) GetUserByEmail(email string) (models.User, error) {
	collection := db.GetCollection()

	var filter interface{}
	response := models.User{}
	filter = bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// User does not exist
			return response, nil
		}
		log.Errorln("Error finding user:", err)
		return response, err
	}

	// Decode the user document into response
	err := result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding user:", err)
		return response, err
	}

	return response, nil
}

func (ur UserRepository) GetUserByID(id string) (models.User, error) {
	collection := db.GetCollection()

	var filter interface{}
	response := models.User{}
	objID, err := primitive.ObjectIDFromHex(id)
	filter = bson.M{"_id": objID}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// User does not exist
			return response, err
		}
		log.Errorln("Error finding user:", err)
		return response, err
	}

	// Decode the user document into response
	err = result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding user:", err)
		return response, err
	}

	return response, nil
}
