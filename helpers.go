package main

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type boshConfig struct {
	Environments []boshEnvironment `yaml:"environments"`
}

type boshEnvironment struct {
	URL          string `yaml:"url"`
	CACert       string `yaml:"ca_cert"`
	Alias        string `yaml:"alias"`
	Username     string `yaml:"password"`
	RefreshToken string `yaml:"refresh_token"`
}

func getBoshConfig() (*boshConfig, error) {
	confFile, err := os.Open(fmt.Sprintf("%s/.bosh/config", os.Getenv("HOME")))
	if err != nil {
		return nil, err
	}

	ret := &boshConfig{}

	err = yaml.NewDecoder(confFile).Decode(ret)
	return ret, err
}
