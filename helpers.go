package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var boshClient *client

type boshConfig struct {
	Environments []boshEnvironment `yaml:"environments"`
}

type boshEnvironment struct {
	URL          string `yaml:"url"`
	CACert       string `yaml:"ca_cert"`
	Alias        string `yaml:"alias"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
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

func getBoshClient(ctx compContext) (*client, error) {
	if boshClient != nil {
		return boshClient, nil
	}

	var envs []string
	var depFound bool
	if envs, depFound = ctx.Flags["--environment"]; !depFound {
		return nil, fmt.Errorf("env not given")
	}
	envName := envs[0]
	cfg, err := getBoshConfig(ctx)
	if err != nil {
		return nil, err
	}

	//I think bosh looks for the address in the alias, and then rescans for the
	// first instance of that address
	// So... first, we look for the alias
	var env *boshEnvironment
	for _, e := range cfg.Environments {
		if e.Alias == envName {
			env = &e
			break
		}
	}

	envAddr := envName
	if env != nil {
		envAddr = env.URL
	}

	log.Write("making client for addr: %s", envAddr)

	ret := &client{
		URL:               envAddr,
		SkipSSLValidation: true,
		cache:             map[string]string{},
	}

	env = nil
	for _, e := range cfg.Environments {
		if e.URL == envAddr {
			env = &e
			break
		}
	}

	if env == nil {
		return nil, fmt.Errorf("Could not get auth info for env: %s", envName)
	}

	ret.Username = env.Username
	ret.Password = env.Password
	ret.RefreshToken = env.RefreshToken

	//--client and --client-secret flags override config
	if client, found := ctx.Flags["--client"]; found {
		ret.Username = client[0]
	}

	if clientSecret, found := ctx.Flags["--client-secret"]; found {
		ret.Password = clientSecret[0]
	}

	boshClient = ret

	return boshClient, nil
}

type boshInstance struct {
	AgentID   string `json:"agent_id"`
	CID       string `json:"cid"`
	Job       string `json:"job"`
	Index     int    `json:"index"`
	ID        string `json:"id"`
	ExpectsVM bool   `json:"expects_vm"`
}

func fetchInstances(c *client, ctx compContext) ([]boshInstance, error) {
	var deployments []string
	var depGiven bool
	if deployments, depGiven = ctx.Flags["--deployment"]; !depGiven {
		return nil, fmt.Errorf("No deployment given")
	}

	ret := []boshInstance{}

	err := c.Get(fmt.Sprintf("/deployments/%s/instances", deployments[0]), &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type boshRelease struct {
	Name     string `json:"name"`
	Versions []struct {
		Version           string `json:"version"`
		CurrentlyDeployed bool   `json:"currently_deployed"`
	} `json:"release_versions"`
}

func fetchReleases(c *client) ([]boshRelease, error) {
	var releases []boshRelease
	err := c.Get(fmt.Sprintf("/releases"), &releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

type boshStemcell struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Deployments []struct {
		Name string `json:"name"`
	} `json:"deployments"`
}

func fetchStemcells(c *client) ([]boshStemcell, error) {
	var stemcells []boshStemcell
	err := c.Get(fmt.Sprintf("/stemcells"), &stemcells)
	if err != nil {
		return nil, err
	}

	return stemcells, nil
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
