package conf

type Config struct {
	HTTP      HTTPConfig
	Upstreams UpstreamsConfig
}

type HTTPConfig struct {
	Addr string
}

type UpstreamsConfig struct {
	Hello   string
	Goodbye string
	Worker  string
}
