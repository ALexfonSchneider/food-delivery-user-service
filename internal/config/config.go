package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"net"
	"net/url"
	"os"
)

type PostgresConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     string `koanf:"port" validate:"required"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
	PoolSize int    `koanf:"poolSize" default:"100"`
}

func (pc *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		pc.User, url.PathEscape(pc.Password),
		net.JoinHostPort(pc.Host, pc.Port),
		pc.Database,
	)
}

type AppEnv string

const (
	Local      AppEnv = "local"
	Production AppEnv = "production"
)

type Config struct {
	appEnv   AppEnv
	Postgres PostgresConfig `koanf:"postgres"`
}

func MustConfig() *Config {
	appEnv := AppEnv(os.Getenv("APP_ENV"))

	if appEnv == "" {
		appEnv = Local
	}

	cfg := &Config{
		appEnv: appEnv,
	}

	k := koanf.New(".")

	path := fmt.Sprintf("config/%s.yaml", cfg.appEnv)

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		panic(err)
	}

	cfg.appEnv = AppEnv(os.Getenv("APP_ENV"))

	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.Struct(cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}
