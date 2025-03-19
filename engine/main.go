package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"modb"

	"sys"
	"web"

	log "github.com/tengfei-xy/go-log"
	"gopkg.in/yaml.v3"
)

var app sys.App

func init_flag() {
	flag.IntVar(&app.Loglevel, "v", log.LEVELINFOINT, fmt.Sprintf("日志等级,%d-%d", log.LEVELFATALINT, log.LEVELDEBUG3INT))
	flag.StringVar(&app.Configpath, "c", "config.yaml", "配置文件路径")
	flag.Parse()
}
func init_config() {
	// 读取配置文件
	f, err := os.ReadFile(app.Configpath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(f, &app.Config)
	if err != nil {
		log.Fatalf("读取配置文件失败:%s", err)
	}
	if app.Config.Web.Server.Port <= 0 {
		app.Config.Web.Server.Port = 21520
	}
	app.Config.Web.SetFullAddress()

}
func init_log() {
	log.SetLevelInt(app.Loglevel)
	_, g := log.GetLevel()
	fmt.Printf("日志等级:%s\n", g)
}
func init_mongo() {
	log.Infof("mongo连接中...")
	str := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		app.Config.DB.User,
		app.Config.DB.Password,
		app.Config.DB.Address,
		app.Config.DB.Port,
		app.Config.DB.Database,
	)
	err := modb.Init(str)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("mongo连接成功!!")
}

func init_web() {

	log.Infof("API: %s", app.Config.Web.Server.FullAddress)
	log.Info("启动监听...")
	web.Init(app.Config.Web)
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
