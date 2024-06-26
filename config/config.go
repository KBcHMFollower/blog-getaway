package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string       `yaml:"env" env-default:"local"`
	HttpServer   HttpServer   `yaml:"http_server"`
	GrpcServices GrpcServices `yaml:"grpc_services"`
}

type HttpServer struct {
	PublicUrl    string        `yaml:"public_url"`
	Address      string        `yaml:"address" env-default:"localhost:8082"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type GrpcServices struct {
	UserService GrpcService `yaml:"user_service"`
	PostService GrpcService `yaml:"post_service"`
}

type GrpcService struct {
	Addr string `yaml:"addr" env-required:"true"`
}

func MustLoad(configPath string) *Config {
	var cfg Config
	if configPath == "" {
		panic("config file path is empty")
	}
	fmt.Println(configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist")
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config")
	}

	return &cfg
}
