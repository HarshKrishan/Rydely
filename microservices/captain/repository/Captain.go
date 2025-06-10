package repository

import (
	"captain/db"
	"captain/models"
	"context"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CaptainRepository struct{}

func (ur CaptainRepository) CheckCaptainExists(email string) (bool, error) {

	collection := db.GetCollection()

	var filter interface{}
	response := models.Captain{}
	filter = bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// Captain does not exist
			return false, nil
		}
		log.Errorln("Error finding Captain:", err)
		return false, err
	}

	// Decode the Captain document into response
	err := result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding Captain:", err)
		return false, err
	}

	return true, nil
}

func (ur CaptainRepository) CreateCaptain(Captain models.Captain) error {
	collection := db.GetCollection()

	// Insert the Captain document into the collection
	_, err := collection.InsertOne(context.TODO(), Captain)
	if err != nil {
		log.Errorln("Error inserting Captain:", err)
		return err
	}

	log.Infoln("Captain created successfully")
	return nil
}

func (ur CaptainRepository) GetCaptainByEmail(email string) (models.Captain, error) {
	collection := db.GetCollection()

	var filter interface{}
	response := models.Captain{}
	filter = bson.M{"email": email}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// Captain does not exist
			return response, nil
		}
		log.Errorln("Error finding Captain:", err)
		return response, err
	}

	// Decode the Captain document into response
	err := result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding Captain:", err)
		return response, err
	}

	return response, nil
}

func (ur CaptainRepository) GetCaptainByID(id string) (models.Captain, error) {
	collection := db.GetCollection()

	var filter interface{}
	response := models.Captain{}
	objID, err := primitive.ObjectIDFromHex(id)
	filter = bson.M{"_id": objID}

	result := collection.FindOne(context.TODO(), filter)
	log.Println("Finding Captain by ID:", id)
	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// Captain does not exist
			log.Println("Captain not found with ID:", id)
			return response, err
		}
		log.Errorln("Error finding Captain:", err)
		return response, err
	}

	// Decode the Captain document into response
	err = result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding Captain:", err)
		return response, err
	}
	log.Println("Captain found:", response)
	return response, nil
}

func (ur CaptainRepository) UpdateCaptain(captain models.Captain) error {
	collection := db.GetCollection()

	objID, err := primitive.ObjectIDFromHex(captain.ID)
	filter := bson.M{"_id": objID}
	// update := bson.M{"$set": captain}
	update := bson.M{
		"$set": bson.M{
			"isactive": captain.IsActive,
		},
	}

	// Update the Captain document in the collection
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorln("Error updating Captain:", err)
		return err
	}

	log.Infoln("Captain updated successfully")
	return nil
}
