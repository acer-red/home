package app

import "github.com/acer-red/official/engine/service/web"

type App struct {
	Loglevel   int
	Configpath string
	Config     Config
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Web web.Config `yaml:"web"`

	DB struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`

	Redis Redis `yaml:"redis"`
}
