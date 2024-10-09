package lib

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var Cfg *ini.File

type ServerConfig struct {
	RunMode        string
	UploadLocation string
	HttpPort       int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	DBType         string
	User           string
	Password       string
	Host           string
	DBName         string
	TablePrefix    string
	RedisHost      string
	RedisIndex     string
}

var Config ServerConfig

func LoadServerConfig() ServerConfig {
	var err error

	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatal(2, "Fail to paser conf/app.ini: %v", err)
	}
	app, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatal(2, "Fail to get section app: %v", err)
	}
	server, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatal(2, "Fail to get section server: %v", err)
	}
	database, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section database: %v", err)
	}
	redis, err := Cfg.GetSection("redis")
	if err != nil {
		log.Fatal(2, "Fail to get section redis: %v", err)
	}

	Config = ServerConfig{
		RunMode:        Cfg.Section("").Key("RUN_MODE").MustString("debug"),
		UploadLocation: app.Key("LOCATION").MustString(""),
		HttpPort:       server.Key("HTTP_PORT").MustInt(),
		ReadTimeout:    time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second,
		WriteTimeout:   time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second,
		DBType:         database.Key("TYPE").MustString(""),
		User:           database.Key("USER").MustString(""),
		Password:       database.Key("PASSWORD").MustString(""),
		Host:           database.Key("HOST").MustString(""),
		DBName:         database.Key("NAME").MustString(""),
		TablePrefix:    database.Key("TABLE_PREFIX").MustString(""),
		RedisHost:      redis.Key("HOST").MustString(""),
		RedisIndex:     redis.Key("INDEX").MustString(""),
	}

	return Config
}
