package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func getBoshConfig(ctx compContext) (*boshConfig, error) {
	location := fmt.Sprintf("%s/.bosh/config", os.Getenv("HOME"))
	if cfg, found := ctx.Flags["--config"]; found {
		location = cfg[0]
	}

	confFile, err := os.Open(location)
	if err != nil {
		return nil, err
	}

	ret := &boshConfig{}

	err = yaml.NewDecoder(confFile).Decode(ret)
	return ret, err
}

func substituteHomeDir(cur string) string {
	if strings.HasPrefix(cur, "~") {
		homeDir := os.Getenv("HOME")
		if !strings.HasSuffix(homeDir, "/") {
			homeDir = homeDir + "/"
		}

		cur = strings.Replace(cur, "~", homeDir, 1)
	}
	return cur
}

func walkDirs(cur string, showLeaf bool, chooseDir bool) ([]string, error) {
	dontFilterPrefix = true
	cur = substituteHomeDir(cur)

	rawPathParts := strings.Split(cur, "/")
	pathParts := []string{}
	for i := 0; i < len(rawPathParts)-1; i++ {
		if rawPathParts[i] != "" {
			pathParts = append(pathParts, rawPathParts[i])
		}
	}
	pathParts = append(pathParts, rawPathParts[len(rawPathParts)-1])
	prefix := ""
	if len(pathParts) > 0 {
		prefix = strings.Join(pathParts[:len(pathParts)-1], "/")
	}
	toList := prefix
	if prefix == "" {
		toList = "./"
	} else {
		prefix += "/"
	}

	cur = prefix + pathParts[len(pathParts)-1]

	log.Write("Prefix: %s", prefix)

	log.Write("Listing %s", toList)

	contents, err := ioutil.ReadDir(toList)
	if err != nil {
		return nil, err
	}

	candidates := []string{prefix + "./", prefix + "../"}
	if chooseDir && prefix != "" {
		candidates = append(candidates, prefix)
	}
	for _, content := range contents {
		slash := ""
		if content.IsDir() {
			slash = "/"
		}

		toAdd := fmt.Sprintf("%s%s%s", prefix, content.Name(), slash)
		candidates = append(candidates, toAdd)
	}

	ret := []string{}
	for _, val := range candidates {
		if (showLeaf || strings.HasSuffix(val, "/")) && strings.HasPrefix(val, cur) {
			ret = append(ret, val)
		}
	}

	if len(ret) == 1 && strings.HasSuffix(ret[0], "/") {
		if !chooseDir {
			ret, err = walkDirs(ret[0], showLeaf, chooseDir)
			if err != nil {
				return nil, err
			}
		} else {
			dontAddSpace = true
		}
	}

	return ret, nil
}
