package main

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/answersuck/internal/app"
	"github.com/ysomad/answersuck/internal/config"
)

func main() {
	var flags app.Flags
	flag.BoolVar(&flags.InDocker, "docker", false, "enter true if app is running inside docker")

	flag.Parse()

	var conf config.Config
	if err := cleanenv.ReadConfig("./configs/local.yml", &conf); err != nil {
		log.Fatalf("config parse error: %s", err)
	}

	app.Run(&conf, flags)
}
