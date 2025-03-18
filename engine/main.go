package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"modb"

	"web"

	log "github.com/tengfei-xy/go-log"
	"gopkg.in/yaml.v3"
)

type App struct {
	loglevel   int
	configpath string
	config     Config
}
type Config struct {
	WS struct {
		Address     string `yaml:"address"`
		SslEnable   bool   `yaml:"ssl_enable"`
		CrtFile     string `yaml:"crt_file"`
		KeyFile     string `yaml:"key_file"`
		Port        int    `yaml:"port"`
		fullAddress string
	} `yaml:"webserver"`
	WC struct {
		AllowOrigin string `yaml:"allow_origin"`
	} `yaml:"webcors"`
	DB struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

var app App

func init_flag() {
	flag.IntVar(&app.loglevel, "v", log.LEVELINFOINT, fmt.Sprintf("日志等级,%d-%d", log.LEVELFATALINT, log.LEVELDEBUG3INT))
	flag.StringVar(&app.configpath, "c", "config.yaml", "配置文件路径")
	flag.Parse()
}
func init_config() {
	// 读取配置文件
	f, err := os.ReadFile(app.configpath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, &app.config)
	if err != nil {
		log.Fatalf("读取配置文件失败:%s", err)
	}
	if app.config.WS.Port <= 0 {
		app.config.WS.Port = 21520
	}

	if app.config.WS.SslEnable {
		app.config.WS.fullAddress = fmt.Sprintf("https://%s:%d", app.config.WS.Address, app.config.WS.Port)
	} else {
		app.config.WS.fullAddress = fmt.Sprintf("http://%s:%d", app.config.WS.Address, app.config.WS.Port)
	}

}
func init_log() {
	log.SetLevelInt(app.loglevel)
	_, g := log.GetLevel()
	fmt.Printf("日志等级:%s\n", g)
}
func init_mongo() {
	log.Infof("mongo连接中...")
	str := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		app.config.DB.User,
		app.config.DB.Password,
		app.config.DB.Address,
		app.config.DB.Port,
		app.config.DB.Database,
	)
	err := modb.Init(str)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("mongo连接成功!!")
}

func init_web() {

	log.Infof("API: %s", app.config.WS.fullAddress)
	log.Info("启动监听...")
	log.Debug3f("cors allow origin: %s", app.config.WC.AllowOrigin)

	web.Init(web.Env{
		CORSAllowOrigin:   app.config.WC.AllowOrigin,
		FullServerAddress: app.config.WS.fullAddress,
		SslEnable:         app.config.WS.SslEnable,
		CrtFile:           app.config.WS.CrtFile,
		KeyFile:           app.config.WS.KeyFile,
		Port:              app.config.WS.Port,
	})
}
func quit() {
	// 创建一个通道来接收信号通知
	sigs := make(chan os.Signal, 1)

	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGINT)
	log.Infof("PID: %d", os.Getpid())
	// 阻塞等待信号
	sig := <-sigs
	fmt.Println(sig)

	err := modb.Disconnect()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(1)
}

func main() {
	init_flag()
	init_config()
	init_log()
	init_mongo()
	go quit()
	init_web()

}
