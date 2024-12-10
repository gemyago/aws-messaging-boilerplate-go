package config

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func LoadTestConfig() *viper.Viper {
	cfg := New()
	lo.Must0(Load(cfg, NewLoadOpts().WithEnv("test")))
	return cfg
}
