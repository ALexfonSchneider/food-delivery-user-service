package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
)

type AppEnv string

const (
	Local      AppEnv = "local"
	Docker     AppEnv = "docker"
	Production AppEnv = "production"
)

type Config struct {
	appEnv        AppEnv
	Postgres      PostgresConfig `koanf:"postgres"`
	Auth          AuthConfig     `koanf:"auth"`
	OpenTelemetry OTELConfig     `koanf:"otel"`
	App           AppConfig      `koanf:"app"`
}

func MustConfig() *Config {
	selectedEnv := Local
	if os.Getenv("APP_ENV") != "" {
		selectedEnv = AppEnv(os.Getenv("APP_ENV"))
	}

	cfg := &Config{
		appEnv: selectedEnv,
	}

	k := koanf.New(".")

	configPath := fmt.Sprintf("config/%s.yaml", cfg.appEnv)
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		panic(err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(cfg); err != nil {
		panic(err)
	}

	fmt.Println(cfg)

	return cfg
}
