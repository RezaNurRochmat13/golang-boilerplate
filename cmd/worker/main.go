package main

import (
	"golang-boilerplate-example/database"
	"golang-boilerplate-example/module/email"
	"log"

	"github.com/hibiken/asynq"
)

func main() {
	server := asynq.NewServer(
		database.RedisOpt(),
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(email.TypeSendEmail, email.HandleSendEmail)

	log.Println("Email worker running...")
	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}
}
