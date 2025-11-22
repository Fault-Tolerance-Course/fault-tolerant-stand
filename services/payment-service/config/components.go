package config

import "time"

type Graceful struct {
	Timeout time.Duration `yaml:"timeout"`
}

type GrpcServer struct {
	Host              string        `yaml:"host"`
	Port              int           `yaml:"port" env:"GRPC_PORT"`
	MaxConnectionIdle time.Duration `yaml:"max_connection_idle"`
	MaxConnectionAge  time.Duration `yaml:"max_connection_age"`
	Timeout           time.Duration `yaml:"timeout"`
	Time              time.Duration `yaml:"time"`
}

type HttpServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port" env:"HTTP_PORT"`
}

type Kafka struct {
	Brokers       []string `yaml:"brokers"`
	ConsumerGroup string   `yaml:"consumer_group"`
}

type Postgres struct {
	Host               string        `yaml:"host"`
	Port               int32         `yaml:"port"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Database           string        `yaml:"database"`
	MaxConnections     int32         `yaml:"max_connections"`
	MinConnections     int32         `yaml:"min_connections"`
	MaxIdleConnections int32         `yaml:"max_idle_connections"`
	MaxConnLifetime    time.Duration `yaml:"max_conn_lifetime"`
}
