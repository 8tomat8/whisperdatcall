package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func main() {
	app := fiber.New()

	cfg := config{}
	envconfig.MustProcess("", cfg)
	m := magic{
		BananaAPIKey:   cfg.BananaAPIKey,
		BananaModelKey: cfg.BananaModelKey,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app.Post("/transcript/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return errors.New("dude, your id sucks")
		}
		body := c.Request().Body()

		text, err := m.transcribe(body)
		if err != nil {
			return errors.Wrap(err, "transcribe")
		}
		return c.SendString(text)
	})

	app.Get("/transcript/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return errors.New("dude, your id sucks")
		}

		return c.SendString("I'm not implemented")
	})

	app.Get("/transcript", func(c *fiber.Ctx) error {
		return c.SendString("I'm not implemented")
	})

	app.Listen(":3000")
}

type magic struct {
	BananaAPIKey   string `envconfig:"BANANA_API_KEY"`
	BananaModelKey string `envconfig:"BANANA_MODEL_KEY"`
}

func (h magic) transcribe(audio []byte) (string, error) {
	// TODO: implement
	return "", nil
}
