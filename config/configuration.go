package config

var Config = new(Configuration)

type Configuration struct {
	Application Application `mapstructure:"application" json:"app" yaml:"app"`
	Log         Log         `mapstructure:"log" json:"log" yaml:"log"`
	Database    Database    `mapstructure:"database" json:"database" yaml:"database"`
	Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	Jwt         Jwt         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
