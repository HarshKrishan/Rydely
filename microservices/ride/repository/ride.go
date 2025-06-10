package repository

import (
	"context"
	"ride/db"
	"ride/models"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RideRepository struct{}

func (ur RideRepository) CheckRideExists(id string) (bool, error) {

	collection := db.GetCollection()

	var filter interface{}
	response := models.Ride{}
	objID, err := primitive.ObjectIDFromHex(id)
	filter = bson.M{"_id": objID}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// Ride does not exist
			return false, nil
		}
		log.Errorln("Error finding ride:", err)
		return false, err
	}

	// Decode the ride document into response
	err = result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding ride:", err)
		return false, err
	}

	return true, nil
}

func (ur RideRepository) CreateRide(ride *models.Ride) error {
	collection := db.GetCollection()

	// Insert the ride document into the collection
	result, err := collection.InsertOne(context.TODO(), ride)
	if err != nil {
		log.Errorln("Error inserting ride:", err)
		return err
	}
	ride.ID = result.InsertedID.(primitive.ObjectID).Hex()

	log.Infoln("Ride created successfully")
	return nil
}

func (ur RideRepository) GetRideByID(rideId string) (models.Ride, error) {
	collection := db.GetCollection()
	var filter interface{}
	response := models.Ride{}
	log.Println("Get ride request for id: ", rideId)
	objID, err := primitive.ObjectIDFromHex(rideId)
	filter = bson.M{"_id": objID}

	result := collection.FindOne(context.TODO(), filter)

	// Check if no document found
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			// ride does not exist
			return response, err
		}
		log.Errorln("Error finding ride:", err)
		return response, err
	}

	// Decode the ride document into response
	err = result.Decode(&response)
	if err != nil {
		log.Errorln("Error decoding ride:", err)
		return response, err
	}
	log.Println("ride found: ", response)
	return response, nil
}

func (ur RideRepository) UpdateRide(ride models.Ride) error {
	collection := db.GetCollection()

	objID, err := primitive.ObjectIDFromHex(ride.ID)
	filter := bson.M{"_id": objID}
	// update := bson.M{"$set": ride}
	update := bson.M{
		"$set": bson.M{
			"status":     ride.Status,
			"captain_id": ride.CaptainID,
		},
	}

	// Update the Captain document in the collection
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorln("Error updating Ride:", err)
		return err
	}

	log.Infoln("Ride updated successfully")
	return nil
}
