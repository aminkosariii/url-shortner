package routes

import (
	"github.com/aminkosariii/url-shortner/database"
	"github.com/aminkosariii/url-shortner/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Request struct {
	URL         string        `json:"url"`
	Customshort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL            string        `json:"URL"`
	Customshort    string        `json:"short"`
	RateRemaining  int           `json:"rate_limit"`
	Expiry         time.Duration `json:"expiry"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx) error {
	var body Request

	//Parse the given request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "can not parse body"})
		//log.Fatalf("can not parse body")
		//return c.Status(fiber.StatusBadRequest)
	}
	//implement rate limiting
	r2 := database.CreateDBClient(1)
	defer func(r2 *redis.Client) {
		err := r2.Close()
		if err != nil {
			log.Fatal("can not close r2")
		}
	}(r2)
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	// If Quota has been exhausted
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("APP_QUOTA"), 30*60*time.Second).Err()
	} else {
		//val, _ := r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":           "rate limit exceeded",
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}
	//check if the input is actual URl
	if _, err := url.ParseRequestURI(body.URL); err != nil {
		panic(err)
	}
	//if !govalidator.IsURL(body.URL) {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	//}

	//check the domain server to avoid infinite loop
	if !helpers.RemoveDomainError(body.URL) {
		//panic("You can't hack my service!!")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You can not hack my service"})
	}
	//enforce https,SSL
	helpers.EnforceHttp(body.URL)
	var id string
	if body.Customshort == "" {

		id = uuid.New().String()[:6]
	} else {
		id = body.Customshort
	}
	r := database.CreateDBClient(0)
	r.Close()

	val1, _ := r.Get(database.Ctx, id).Result()
	if val1 != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "this short url already used"})

	}
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	resp := Response{
		URL:            body.URL,
		Expiry:         body.Expiry,
		Customshort:    "",
		RateRemaining:  10,
		RateLimitReset: 30,
	}
	r2.Decr(database.Ctx, c.IP())
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.RateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.RateLimitReset = ttl / time.Nanosecond / time.Minute
	return c.Status(fiber.StatusOK).JSON(resp)

}
