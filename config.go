package evediscordkm

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Channels   []ChannelConfig  `yaml:"channels"`
	ZKillboard ZKillboardConfig `yaml:"zkillboard"`
}

type ZKillboardConfig struct {
	ReplayFrom string `yaml:"replay_from"`
}

type ChannelConfig struct {
	Type        string
	Constraints Constraints
	Config      map[string]interface{} `json:"-"`
}

func ReadFile(fileName string) Config {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	return config
}
