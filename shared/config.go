package shared

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}

type Database struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}

type Server struct {
	Address string `mapstructure:"address"`
}

func LoadConfig(paths ...string) (Config, error) {
	if len(paths) == 0 {
		return Config{}, fmt.Errorf("indique ao menos um caminho de configuração")
	}

	var config Config
	v := viper.New()
	v.SetConfigName("gestao")
	// v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, p := range paths {
		v.AddConfigPath(os.ExpandEnv(p))
	}

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("falha ao ler configuração gestao: %w", err)
	}
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("falha ao interpretar configuração gestao: %w", err)
	}

	return config, nil
}
