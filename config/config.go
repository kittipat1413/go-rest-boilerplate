package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	*log.Logger
	Viper *viper.Viper
}

const (
	alwayYesKey         = "ALWAYS_YES"
	listenAddrKey       = "LISTEN_ADDR"
	databaseUrlKey      = "DATABASE_URL"
	apiPrefixKey        = "API_PREFIX"
	backofficeApiKey    = "BACKOFFICE_API_KEY"
	backofficeApiSecret = "BACKOFFICE_API_SECRET"
)

func MustConfigure() *Config {
	if cfg, err := Configure(); err != nil {
		log.Fatalln(err)
		return nil
	} else {
		return cfg
	}
}

func Configure() (*Config, error) {
	v, err := initViper()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &Config{
		Logger: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		Viper:  v,
	}, nil
}

func initViper() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()
	if _, err := os.Stat("env.yaml"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile("env.yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	} else if _, err := os.Stat("../env.yaml"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile("../env.yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	v.SetDefault(alwayYesKey, "0")
	v.SetDefault(listenAddrKey, ":4000")
	v.SetDefault(databaseUrlKey, "postgres:///go-rest-boilerplate?sslmode=disable")
	v.SetDefault(apiPrefixKey, "http://0.0.0.0:4000")
	return v, nil
}

func (c *Config) AlwaysYes() bool             { return c.Viper.GetBool(alwayYesKey) }
func (c *Config) ListenAddr() string          { return c.Viper.GetString(listenAddrKey) }
func (c *Config) DatabaseURL() string         { return c.Viper.GetString(databaseUrlKey) }
func (c *Config) APIPrefix() string           { return c.Viper.GetString(apiPrefixKey) }
func (c *Config) BackofficeApiKey() string    { return c.Viper.GetString(backofficeApiKey) }
func (c *Config) BackofficeApiSecret() string { return c.Viper.GetString(backofficeApiSecret) }

func (c *Config) AllConfigurations() map[string]interface{} {
	m := map[string]interface{}{}
	for _, key := range c.Viper.AllKeys() {
		m[key] = c.Viper.Get(key)
	}
	return m
}
