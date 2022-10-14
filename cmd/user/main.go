package main

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/answersuck/internal/user/app"
	"github.com/ysomad/answersuck/internal/user/config"
)

func main() {
	var (
		conf     config.Config
		confPath string
	)

	flag.StringVar(
		&confPath,
		"config",
		"config/local.yml",
		"path to yml config file",
	)

	if err := cleanenv.ReadConfig(confPath, &conf); err != nil {
		log.Fatalf("config parse error: %s", err)
	}

	app.Run(&conf)
}
