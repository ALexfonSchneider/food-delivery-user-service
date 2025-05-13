package config

type PostgresConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     string `koanf:"port" validate:"required"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
	PoolSize *int   `koanf:"poolSize"`
}

func (p *PostgresConfig) GetHost() string {
	return p.Host
}

func (p *PostgresConfig) GetPort() string {
	return p.Port
}

func (p *PostgresConfig) GetUser() string {
	return p.User
}

func (p *PostgresConfig) GetPassword() string {
	return p.Password
}

func (p *PostgresConfig) GetDatabase() string {
	return p.Database
}

func (p *PostgresConfig) GetPoolSize() int {
	if p.PoolSize != nil {
		return *p.PoolSize
	}
	return 100
}
