package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address string `yaml:"address" env-default:"localhost:8082"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoad() *Config{
	configPath:=fetchConfigPath()
	if configPath == ""{
		panic("config path is empty")
	}

	if _,err := os.Stat(configPath); os.IsNotExist(err){
		panic("config file does not exist")
	}

	var cfg Config
	
	if err := cleanenv.ReadConfig(configPath, &cfg); err !=nil{
		panic("cannot read config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == ""{
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}