package config

type Kubernetes struct {
	IP string `mapstructure:"ip"`
	Token string `mapstructure:"token"`
}
