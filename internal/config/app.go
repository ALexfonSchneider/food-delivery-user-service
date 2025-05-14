package config

type AppConfig struct {
	GRPCPort string `koanf:"grpc_port"`
	GRPCHost string `koanf:"grpc_host"`
}
