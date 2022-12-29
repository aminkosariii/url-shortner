package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
)

var Ctx = context.Background()

func CreateDBClient(dbNumber int) *redis.Client {

	// create connection to database
	Client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNumber,
	})
	//test connection to redis
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatal(http.StatusInternalServerError, err)
	}

	return Client
}
