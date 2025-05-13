package config

type OTELConfig struct {
	GRPCEndpoint string `koanf:"grpc_endpoint" required:"true"`
}
