package server

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"time"
)

type Server struct {
	RunMode      string
	Debug        bool
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TrustProxies []string
}

type Database struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type App struct {
	AppTitle         string
	AppDeveloper     string
	AppDeveloperMail string
}

var Config = &Server{}
var DatabaseConfig = &Database{}
var AppConfig = &App{}

var cfg *ini.File

func Init() {
	cfg = loadAppConfigIni()

	mapTo("server", Config)
	mapTo("database", DatabaseConfig)
	mapTo("app", AppConfig)

	configToGoReadable()
}

func loadAppConfigIni() *ini.File {
	path := getConfigFilePath()
	cfg, err := ini.Load(path)
	if err != nil {
		panic("fail to parse: " + path)
	}
	return cfg
}

func getConfigFilePath() string {
	envConfig := os.Getenv("CONFIG")
	if envConfig == "" {
		envConfig = "local"
	}
	return fmt.Sprintf("config/app_%s.ini", envConfig)
}

func configToGoReadable() {
	Config.WriteTimeout = Config.WriteTimeout * time.Second
	Config.ReadTimeout = Config.ReadTimeout * time.Second

	if Config.RunMode == "debug" {
		Config.Debug = true
	} else {
		Config.Debug = false
	}
}

func mapTo(section string, config interface{}) {
	err := cfg.Section(section).MapTo(config)
	if err != nil {
		panic("fail to load section: " + section)
	}
}
