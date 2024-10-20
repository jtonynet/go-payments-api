package config

import (
	"github.com/spf13/viper"
)

type API struct {
	Name       string `mapstructure:"API_NAME"`
	Port       string `mapstructure:"API_PORT"`
	TagVersion string `mapstructure:"API_TAG_VERSION"`
	Env        string `mapstructure:"ENV"`
}

type Database struct {
	Strategy string `mapstructure:"DATABASE_STRATEGY"`
	Driver   string `mapstructure:"DATABASE_DRIVER"`

	Host    string `mapstructure:"DATABASE_HOST"`
	User    string `mapstructure:"DATABASE_USER"`
	Pass    string `mapstructure:"DATABASE_PASSWORD"`
	DB      string `mapstructure:"DATABASE_DB"`
	Port    string `mapstructure:"DATABASE_PORT"`
	SSLmode string `mapstructure:"DATABASE_SSLMODE"`
}

type Config struct {
	API      API      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	env := viper.GetString("ENV")
	switch env {
	case "test":
		viper.SetConfigName(".env.TEST")
	case "dev", "":
		viper.SetConfigName(".env")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
