package conf

type Config struct {
    HTTP HTTPConfig
    GRPC GRPCConfig
}

type HTTPConfig struct {
    Addr string // ":8081"
}

type GRPCConfig struct {
    Addr string // ":9081"
}
