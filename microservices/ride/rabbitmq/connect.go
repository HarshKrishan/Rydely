package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"ride/models"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

func Init() {
	var err error
	rabbitmq_uri := os.Getenv("RABBITMQ_URI")
	Conn, err = amqp.Dial(rabbitmq_uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	log.Println("Connected to RabbitMQ")
}

func PublishRide(ride models.Ride, queue_name string) error {
	q, err := Channel.QueueDeclare(
		queue_name, // queue name
		true,         // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(ride)
	if err != nil {
		return err
	}

	err = Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Println("Ride published to queue")
	return nil
}

func ConsumeRideQueue() {
	q, err := Channel.QueueDeclare(
		"ride_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Queue declare failed: %v", err)
	}

	msgs, err := Channel.Consume(
		q.Name,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	go func() {
		for msg := range msgs {
			var ride models.Ride
			if err := json.Unmarshal(msg.Body, &ride); err != nil {
				log.Println("Failed to decode ride:", err)
				continue
			}

			log.Printf("Received new ride: %+v", ride)

		}
	}()
}
