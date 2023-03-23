package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/exp/slog"
)

const OpenAIAPIKey = "lol"

func main() {
	app := fiber.New()

	cfg := &config{}
	envconfig.MustProcess("", cfg)
	m := magic{
		BananaAPIKey:   cfg.BananaAPIKey,
		BananaModelKey: cfg.BananaModelKey,
		OpenAICli:      openai.NewClient(OpenAIAPIKey),
	}
	convID := "f8086bfa-959a-4b8c-8fa6-101c224ff2bc"
	summary, err := m.summarize(context.Background(), someText)
	if err != nil {
		log.Fatal(errors.Wrap(err, "summarize"))
	}
	slog.Info("finished summarize", "summary", summary)

	if err = sendDataToConnector(convID, summary); err != nil {
		slog.Error("send data to the connector", err)
	}
	slog.Info("Success!")
	return
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

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

func sendDataToConnector(convID, text string) error {
	const uri = "nope"

	type req struct {
		ConvID  string `json:"conversationId"`
		Summary string `json:"summary"`
	}

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(req{ConvID: convID, Summary: text}); err != nil {
		return errors.Wrap(err, "encode the request")
	}

	resp, err := http.Post(uri, "application/json", buf)
	if err != nil {
		return errors.Wrap(err, "call the webhook")
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Printf("Failed webhook request with body:\n%s\n", string(body))
		}
		return errors.Errorf("received %d from the webhook", resp.StatusCode)
	}
	return nil
}
