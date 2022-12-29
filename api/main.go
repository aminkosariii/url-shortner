package main

import (
	"github.com/aminkosariii/url-shortner/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveUrl)
	app.Post("/api/v1", routes.ShortenUrl)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can not load envs")
	}
	app := fiber.New()
	// middleware for logging HTTP request/response details and displaying results in a file
	loggeroutputFile, err := os.OpenFile("loggeroutputFile.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error while opening log file: %v", err)
	}
	loggerConfig := logger.Config{
		Output: loggeroutputFile,
	}
	app.Use(logger.New(loggerConfig))
	//Use Limiter middleware to limit repeated requests
	//limiterConfig := limiter.Config{Next: func(c *fiber.Ctx) bool {
	//	return c.IP() == os.Getenv("DOMAIN")
	//},
	//	Max:        10,
	//	Expiration: 30 * time.Second,
	//	LimitReached: func(c *fiber.Ctx) error {
	//		return c.SendFile("./templates/ratelimit.html")
	//	},
	//}
	//app.Use(limiter.New(limiterConfig))

	SetupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
