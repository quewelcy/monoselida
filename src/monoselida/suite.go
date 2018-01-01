package monoselida

import (
	"log"

	"gopkg.in/yaml.v2"
)

//Config describes xpath rule for retrieving data
//Rule xpath rule definition
//Repeat look for single entry or multiple
//IgnoreText get content of tag or not, false for wrapper
//IsTitle rule is for title
//Configs enclosed configs
type Config struct {
	Rule       string
	NextPage   string
	Repeat     bool
	IgnoreText bool
	IsTitle    bool
	Configs    []Config
}

//ProcessSuite unmarshals config yaml into Config
func ProcessSuite(suiteDef []byte) Config {
	var config Config
	err := yaml.Unmarshal(suiteDef, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
