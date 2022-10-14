package main

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/answersuck/user/internal/app"
	"github.com/ysomad/answersuck/user/internal/config"
)

func main() {
	var (
		conf     config.Config
		confPath string
	)

	flag.StringVar(
		&confPath,
		"config",
		"./configs/local.yml",
		"path to yml config file",
	)

	if err := cleanenv.ReadConfig(confPath, &conf); err != nil {
		log.Fatalf("config parse error: %s", err)
	}

	app.Run(&conf)
}
