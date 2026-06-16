package conf

type Config struct {
	HTTP  HTTPConfig
	Redis RedisConfig
}

type HTTPConfig struct {
	Addr string
}

type RedisConfig struct {
	Addr string
}
