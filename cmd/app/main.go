package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/quizlyfun/quizly-backend/internal/app"
	"github.com/quizlyfun/quizly-backend/internal/config"
)

func main() {
	var cfg config.Aggregate

	err := cleanenv.ReadConfig("./configs/local.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&cfg)
}
