package main

import (
	"2024_1_IMAO/internal/app"
	"log"
)

func main() {
	srv := new(app.Server)

	err := srv.Run()

	if err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
