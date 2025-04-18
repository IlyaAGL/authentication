package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/agl/authentication/api"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), DATABASE_URL)

	if err != nil {
		log.Fatalf("Cant connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("DB WAS CONNECTED")
	api.SetupRoutes(conn)
}
