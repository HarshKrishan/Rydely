package main

import (
	"ride/db"
	"ride/rabbitmq"
	"ride/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Starting ride service on port 8080...")

	config := AppConfig{}

	config.LoadEnv()
	db.InitMongoDB()
	config.SetLogLevel()
	rabbitmq.Init()
	server.Init()
}
