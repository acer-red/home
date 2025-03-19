package sys

import "fmt"

type App struct {
	Loglevel   int
	Configpath string
	Config     Config
}
type Web struct {
	Server struct {
		Address     string `yaml:"address"`
		SslEnable   bool   `yaml:"ssl_enable"`
		CrtFile     string `yaml:"crt_file"`
		KeyFile     string `yaml:"key_file"`
		Port        int    `yaml:"port"`
		FullAddress string `yaml:"-"`
	} `yaml:"server"`
	CORS struct {
		Enable      bool   `yaml:"enable"`
		AllowOrigin string `yaml:"allow_origin"`
	} `yaml:"cors"`
}

type Config struct {
	Web Web `yaml:"web"`

	DB struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

func (w *Web) SetFullAddress() {
	if w.Server.SslEnable {
		w.Server.FullAddress = fmt.Sprintf("https://%s:%d", w.Server.Address, w.Server.Port)
	} else {
		w.Server.FullAddress = fmt.Sprintf("http://%s:%d", w.Server.Address, w.Server.Port)
	}
}
