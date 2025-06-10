package main

import (
	"user/db"
	"user/rabbitmq"
	"user/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Starting user service on port 8080...")

	config := AppConfig{}

	config.LoadEnv()
	db.InitMongoDB()
	config.SetLogLevel()
	rabbitmq.Init()
	rabbitmq.ConsumeRideQueue()
	server.Init()
}
