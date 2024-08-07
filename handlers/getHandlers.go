package handlers

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	lua "github.com/yuin/gopher-lua"
)

type Handlers struct {
	Port     int     `yaml:"port" env:"PORT" env-default:"8080"`
	Cert     string  `yaml:"cert" env:"CERT" env-default:""`
	Key      string  `yaml:"key"  env:"KEY" env-default:""`
	Handlers []Route `yaml:"handlers"`
}

type Route struct {
	Route       string `yaml:"route"`
	Content     string `yaml:"content"`
	ContentType string `yaml:"contentType"`
	LuaFunc     *lua.LFunction
}

var Cfg Handlers

func GetHandlers(lua bool) (Handlers, error) {
	if lua {
		fmt.Println("Running with lua config")
		return Cfg, nil
	}

	err := cleanenv.ReadConfig("kserver.yml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		if !lua {
			InitConfig()
		} else {
			fmt.Println("Running with lua config")
		}
	}
	err = cleanenv.ReadConfig("kserver.yaml", &Cfg)
	if err == nil {
		return Cfg, nil
	} else {
		if !lua {
			InitConfig()
		} else {
			fmt.Println("Running with lua config")
		}
	}

	return Handlers{}, fmt.Errorf("unable to find or read config file")
}

func InitConfig() {
	f, err := os.Create("kserver.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write([]byte(defCfg))

	fmt.Println("\033[1mConfig created \033[0m")
	fmt.Println("")
	fmt.Println("Edit \033[33mkserver.yml \033[0mand run again")

	os.Exit(0)
}
