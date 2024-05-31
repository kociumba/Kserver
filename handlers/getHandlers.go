package handlers

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Handlers struct {
	Port     int     `yaml:"port"`
	Handlers []Route `yaml:"handlers"`
}

type Route struct {
	Route       string `yaml:"route"`
	Content     string `yaml:"content"`
	ContentType string `yaml:"contentType"`
}

var Cfg Handlers

const defCfg = `
# yaml-language-server: $schema=https://raw.githubusercontent.com/kociumba/kserver/main/.kserver

port: 8080

handlers:
- route: /
  content: ./index.html
  contentType: text/html

`

func GetHandlers() (Handlers, error) {
	err := cleanenv.ReadConfig("kserver.yml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		InitConfig()
	}
	err = cleanenv.ReadConfig("kserver.yaml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		InitConfig()
	}
	err = cleanenv.ReadConfig("kserver.json", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		InitConfig()
	}

	return Handlers{}, fmt.Errorf("unable to read handlers from config files")
}

func InitConfig() {
	f, err := os.Create("kserver.yml")
	if err != nil {
		panic(err)
	}

	f.Write([]byte(defCfg))

	fmt.Println("\033[1mConfig created \033[0m")
	fmt.Println("")
	fmt.Println("Edit \033[33mkserver.yml \033[0mand run again")

	os.Exit(0)
}
