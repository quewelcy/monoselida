package monoselida

import (
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rule    string
	Repeat  bool
	UseText bool
	Configs []Config
}

func ProcessSuite(suiteDef []byte) Config {
	var config Config
	err := yaml.Unmarshal(suiteDef, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
