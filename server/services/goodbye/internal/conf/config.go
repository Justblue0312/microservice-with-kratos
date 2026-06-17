package conf

type Config struct {
	HTTP  HTTPConfig
	GRPC  GRPCConfig
	NATS  NATSConfig
	Redis RedisConfig
	Asynq AsynqConfig
}

type AsynqConfig struct {
	Addr        string
	Concurrency int
}

type RedisConfig struct {
	Addr string
}

type HTTPConfig struct {
	Addr string
}

type GRPCConfig struct {
	Addr string
}

type NATSConfig struct {
	URL string
}
