package config

import (
	"github.com/spf13/viper"
)

type API struct {
	Name       string `mapstructure:"API_NAME"`
	Port       string `mapstructure:"API_PORT"`
	TagVersion string `mapstructure:"API_TAG_VERSION"`
	Env        string `mapstructure:"ENV"`
	TimeoutSLA int    `mapstructure:"API_TIMEOUT_SLA_IN_MS"`
}

type Router struct {
	Strategy string `mapstructure:"HTTP_ROUTER_STRATEGY"` // gin
}

type Logger struct {
	Strategy  string `mapstructure:"LOG_STRATEGY"`   // slog
	Level     string `mapstructure:"LOG_LEVEL"`      // debug | info | warn | error
	Format    string `mapstructure:"LOG_OPT_FORMAT"` // text | json
	AddSource bool   `mapstructure:"LOG_OPT_ADD_SOURCE_BOOL"`
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

type InMemoryDatabase struct {
	Strategy   string
	Pass       string
	Port       string
	Host       string
	DB         int
	Protocol   int
	Expiration int
}

type InMemoryDatabaseConverter interface {
	ToInMemoryDatabase() (InMemoryDatabase, error)
}
type Cache struct {
	Strategy   string `mapstructure:"CACHE_IN_MEMORY_STRATEGY"`
	Pass       string `mapstructure:"CACHE_IN_MEMORY_PASSWORD"`
	Port       string `mapstructure:"CACHE_IN_MEMORY_PORT"`
	Host       string `mapstructure:"CACHE_IN_MEMORY_HOST"`
	DB         int    `mapstructure:"CACHE_IN_MEMORY_DB"`
	Protocol   int    `mapstructure:"CACHE_IN_MEMORY_PROTOCOL"`
	Expiration int    `mapstructure:"CACHE_IN_MEMORY_EXPIRATION_DEFAULT_IN_MS"`
}

func (c *Cache) ToInMemoryDatabase() (InMemoryDatabase, error) {
	return InMemoryDatabase{
		Strategy:   c.Strategy,
		Pass:       c.Pass,
		Port:       c.Port,
		Host:       c.Host,
		DB:         c.DB,
		Protocol:   c.Protocol,
		Expiration: c.Expiration,
	}, nil
}

type Lock struct {
	Strategy   string `mapstructure:"LOCK_IN_MEMORY_STRATEGY"`
	Pass       string `mapstructure:"LOCK_IN_MEMORY_PASSWORD"`
	Port       string `mapstructure:"LOCK_IN_MEMORY_PORT"`
	Host       string `mapstructure:"LOCK_IN_MEMORY_HOST"`
	DB         int    `mapstructure:"LOCK_IN_MEMORY_DB"`
	Protocol   int    `mapstructure:"LOCK_IN_MEMORY_PROTOCOL"`
	Expiration int    `mapstructure:"LOCK_IN_MEMORY_EXPIRATION_DEFAULT_IN_MS"`
}

func (l *Lock) ToInMemoryDatabase() (InMemoryDatabase, error) {
	return InMemoryDatabase{
		Strategy:   l.Strategy,
		Pass:       l.Pass,
		Port:       l.Port,
		Host:       l.Host,
		DB:         l.DB,
		Protocol:   l.Protocol,
		Expiration: l.Expiration,
	}, nil
}

type Config struct {
	API      API      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
	Router   Router   `mapstructure:",squash"`
	Logger   Logger   `mapstructure:",squash"`
	Cache    Cache    `mapstructure:",squash"`
	Lock     Lock     `mapstructure:",squash"`
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
