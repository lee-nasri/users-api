package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mstoykov/envconfig"
)

type Config struct {
	AppName     string `envconfig:"APP_NAME"`
	AppVersion  string `envconfig:"APP_VERSION"`
	AppEnv      string `envconfig:"APP_ENV"`
	AppPort     string `envconfig:"APP_PORT"`
	HTTPTimeout int64  `envconfig:"HTTP_TIMEOUT_IN_MS"`
	Redis       RedisConfig
	Datadog     datadogConfig
	Metric      metricConfig
}

type RedisConfig struct {
	Host         string `envconfig:"REDIS_HOST"`
	Port         string `envconfig:"REDIS_PORT"`
	Username     string `envconfig:"REDIS_USERNAME"`
	Password     string `envconfig:"REDIS_PASSWORD"`
	KeyPrefix    string `envconfig:"REDIS_KEY_PREFIX"`
	Index        string `envconfig:"REDIS_INDEX"`
	DefaultLimit int    `envconfig:"REDIS_DEFAULT_LIMIT"`
}

type datadogConfig struct {
	Env       string `envconfig:"DD_ENV"`
	Service   string `envconfig:"DD_SERVICE"`
	Component string `envconfig:"DD_COMPONENT"`
	Version   string `envconfig:"DD_VERSION"`
	AgentHost string `envconfig:"DD_AGENT_HOST"`
}

type metricConfig struct {
	Port int `envconfig:"METRIC_PORT"`
}

var config Config

func Init() {
	err := godotenv.Load()
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Printf("🟥 read config error: %v", err)
		} else {
			fmt.Println("use environment from OS")
		}
	}
	err = envconfig.Process("", &config)
	if err != nil {
		fmt.Printf("🟥 parse config error: %v", err)
	}
}

func GetConfig() *Config {
	return &config
}
