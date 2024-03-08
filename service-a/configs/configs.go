package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	WEB_SERVER_PORT             string `mapstructure:"WEB_SERVER_PORT"`
	URL_BASE                    string `mapstructure:"URL_BASE"`
	EXTERNAL_CALL_PORT          string `mapstructure:"EXTERNAL_CALL_PORT"`
	EXTERNAL_CALL_URL           string `mapstructure:"EXTERNAL_CALL_URL"`
	OTEL_SERVICE_NAME           string `mapstructure:"OTEL_SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_ENDPOINT string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	MICROSERVICE_NAME           string `mapstructure:"MICROSERVICE_NAME"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
