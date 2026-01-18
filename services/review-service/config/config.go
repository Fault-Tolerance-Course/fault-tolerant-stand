package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"review/internal/pkg/loadshedding"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	GrpcServer GrpcServer `yaml:"grpc_server"`
	HttpServer HttpServer `yaml:"http_server"`

	LoadShedding loadshedding.Config `yaml:"load_shedding"`

	Graceful Graceful          `yaml:"graceful"`
	Targets  map[string]string `yaml:"service"`
}

func Instance() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yaml", instance); err != nil {
			log.Fatalf("read config error: %s", err.Error())
		}
		if err := cleanenv.ReadEnv(instance); err != nil {
			log.Fatalf("read env error: %s", err.Error())
		}

		if instance.Targets != nil {
			return
		}

		instance.Targets = make(map[string]string)

		for _, env := range os.Environ() {
			if !strings.Contains(env, "_SERVICE_ADDR") {
				continue
			}

			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}

			serviceName := strings.ToLower(strings.Split(parts[1], ":")[0])
			instance.Targets[serviceName] = parts[1]
		}
	})
	return instance
}
