package configs

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
)

const (
	defaultConfigPath = "config.yaml"
)

func pathOrDefault(paths []string) []string {
	if len(paths) == 0 {
		return []string{defaultConfigPath}
	}

	return paths
}

func Load(configPaths ...string) (cfg Config) {
	configPaths = pathOrDefault(configPaths)

	k := koanf.New(".")

	if err := k.Load(file.Provider(configPaths[0]), yaml.Parser()); err != nil {
		log.Fatalf("[CONFIG] Failed to load config file: [%v]", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("[CONFIG] Failed to unmarshal config file: [%v]", err)
	}

	return
}
