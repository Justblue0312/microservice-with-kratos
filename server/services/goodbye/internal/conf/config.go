package conf

type Config struct {
	GRPC GRPCConfig
	NATS NATSConfig
}

type GRPCConfig struct {
	Addr string
}

type NATSConfig struct {
	URL string
}
