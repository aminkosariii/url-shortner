package routes

import (
	"github.com/aminkosariii/url-shortner/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateDBClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			log.Fatalf("can not close database")
		}
	}(r)
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "can not find short in database"})
	}
	//test connection to redis

	//}else if err != nil{
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "can not connect to database"})
	//}

	//redirect url
	rInt := database.CreateDBClient(1)
	defer rInt.Close()
	_ = rInt.Incr(database.Ctx, "counter")
	return c.Redirect(value, 301)
}
