package main

type config struct {
	BananaAPIKey   string `envconfig:"BANANA_API_KEY"`
	BananaModelKey string `envconfig:"BANANA_MODEL_KEY"`
}
