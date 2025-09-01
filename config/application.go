package config

type Application struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}
