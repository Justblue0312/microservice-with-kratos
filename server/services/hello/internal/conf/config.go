package conf

type Config struct {
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	Upstream UpstreamConfig
	NATS     NATSConfig
}

type HTTPConfig struct {
	Addr string
}

type GRPCConfig struct {
	Addr string
}

type UpstreamConfig struct {
	Goodbye string
}

type NATSConfig struct {
	URL string
}
