package shared

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("gestao")
}

type Config struct {
	Database Database `mapstructure:"database"`
}

type Database struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}

func LoadConfig(paths ...string) (config Config, er error) {

	if len(paths) == 0 {
		er = fmt.Errorf("indique um caminho")
		return
	}

	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.Unmarshal(&config)

	return
}
