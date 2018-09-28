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

type filepath struct {
	parts    []string
	absolute bool
	dir      bool
}

func parseFilepath(path string) filepath {
	ret := filepath{}
	rawPathParts := strings.Split(path, "/")
	for i := 0; i < len(rawPathParts)-1; i++ {
		if rawPathParts[i] != "" {
			ret.parts = append(ret.parts, rawPathParts[i])
		}
	}
	ret.absolute = strings.HasPrefix(path, "/")
	ret.dir = len(ret.parts) == 0 || strings.HasSuffix(path, "/")
	return ret
}

func (f filepath) String() string {
	if len(f.parts) == 0 && !f.absolute {
		return "./"
	}

	prefix := ""
	if f.absolute {
		prefix = "/"
	}

	return prefix + strings.Join(f.parts, "/")
}

func walkDirs(cur string, showLeaf bool, chooseDir bool) ([]string, error) {
	dontFilterPrefix = true
	searchPath := parseFilepath(substituteHomeDir(cur)).String()
	lastSlash := strings.LastIndex(cur, "/")
	curDir := ""
	if lastSlash >= 0 {
		curDir = cur[:lastSlash] + "/"
	}

	log.Write("Prefix: %s", curDir)
	log.Write("Listing %s", searchPath)

	contents, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return nil, err
	}

	candidates := []string{curDir + "./", curDir + "../"}
	if chooseDir && curDir != "" {
		candidates = append(candidates, curDir)
	}
	for _, content := range contents {
		slash := ""
		if content.IsDir() {
			slash = "/"
		} else if content.Mode()&os.ModeSymlink > 0 {
			derefSymlink, err := os.Stat(searchPath + content.Name())
			if err != nil {
				return nil, err
			}
			if derefSymlink.IsDir() {
				slash = "/"
			}
		}

		toAdd := fmt.Sprintf("%s%s%s", curDir, content.Name(), slash)
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
