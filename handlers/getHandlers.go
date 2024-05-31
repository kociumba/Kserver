package handlers

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Handlers struct {
	Handlers []Route `yaml:"handlers"`
}

type Route struct {
	Route       string `yaml:"route"`
	Content     string `yaml:"content"`
	ContentType string `yaml:"contentType"`
}

var Cfg Handlers

func GetHandlers() (Handlers, error) {
	err := cleanenv.ReadConfig("kserver.yml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cleanenv.ReadConfig("kserver.yaml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cleanenv.ReadConfig("kserver.json", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		fmt.Println(err)
		os.Exit(1)
	}

	return Handlers{}, fmt.Errorf("unable to read handlers from config files")
}
