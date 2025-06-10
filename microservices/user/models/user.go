package models

type User struct {
	ID       string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


type Ride struct {
	ID          string     `json:"id,omitempty" bson:"_id,omitempty"`
	CaptainID   string     `json:"captain_id" bson:"captain_id"` // reference to Captain
	UserID      string     `json:"user_id" bson:"user_id"`       // reference to User
	Pickup      string     `json:"pickup" bson:"pickup"`
	Destination string     `json:"destination" bson:"destination"`
	Status      RideStatus `json:"status" bson:"status"`
}

type RideStatus string

const (
	RideStatusRequested RideStatus = "requested"
	RideStatusAccepted  RideStatus = "accepted"
	RideStatusStarted   RideStatus = "started"
	RideStatusCompleted RideStatus = "completed"
)