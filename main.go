package main

import (
	"log"
	"os"

	"github.com/equalsgibson/five9-go/five9"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	c := five9.NewService(os.Getenv("BASE_API_URL"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	c.Supervisor().WebSocket().Ping()
}
