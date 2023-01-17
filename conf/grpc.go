package conf

type GrpcConfig struct {
	Port int `json:"port" env:"PORT,notEmpty" validate:"required"`
}
