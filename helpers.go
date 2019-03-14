package main

import (
	"fmt"
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

type filepath struct {
	parts    []string
	absolute bool
	dir      bool
}

func parseFilepath(path string) filepath {
	ret := filepath{}
	rawPathParts := strings.Split(path, "/")
	for i := 0; i < len(rawPathParts); i++ {
		ret.parts = append(ret.parts, rawPathParts[i])
	}

	//Trim last section off of paths ending in "/" (directories)
	if len(ret.parts) > 0 && ret.parts[len(ret.parts)-1] == "" {
		ret.dir = true
		ret.parts = ret.parts[:len(ret.parts)-1]
	}

	ret.absolute = strings.HasPrefix(path, "/")
	//Trim unnecessary first part if absolute
	if ret.absolute {
		ret.parts = ret.parts[1:]
	}

	return ret
}

func (f filepath) SearchString() string {
	if len(f.parts) == 0 && !f.absolute {
		return "."
	}

	if len(f.parts) > 0 && f.parts[0] == "~" {
		f.parts = f.parts[1:]
		homeFilepath := parseFilepath(os.Getenv("HOME"))
		f.parts = append(homeFilepath.parts, f.parts...)
		f.absolute = true
	}

	prefix := ""
	if f.absolute {
		prefix = "/"
	}

	suffix := ""
	if !(f.absolute && len(f.parts) == 0) && f.dir {
		suffix = "/"
	}

	return prefix + strings.Join(f.parts, "/") + suffix
}

func (f filepath) OriginalString() string {
	prefix := ""
	if f.absolute {
		prefix = "/"
	}

	suffix := ""
	if !(f.absolute && len(f.parts) == 0) && f.dir {
		suffix = "/"
	}

	return prefix + strings.Join(f.parts, "/") + suffix
}

func (f filepath) GetContents(acceptFiles bool) ([]filepath, error) {
	file, err := os.Open(f.SearchString())
	if err != nil {
		return nil, err
	}

	defer file.Close()

	infos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	ret := []filepath{}
	for _, info := range infos {
		log.Write("INFO NAME: %+v\n", info.Name())
		dir := info.IsDir()
		if info.Mode() == os.ModeSymlink {
			symlinkInfo, err := os.Stat(filepath{
				parts:    append(f.parts, info.Name()),
				absolute: f.absolute,
			}.SearchString())
			if err != nil {
				return nil, err
			}

			dir = symlinkInfo.IsDir()
		}

		if !acceptFiles && !dir {
			continue
		}

		theseParts := make([]string, len(f.parts))
		copy(theseParts, f.parts)
		ret = append(ret, filepath{
			parts:    append(theseParts, info.Name()),
			absolute: f.absolute,
			dir:      dir,
		})
	}
	return ret, nil
}

func walkDirs(cur string, acceptFile bool) ([]string, error) {
	//We're handling our own space additions
	dontAddSpace = true
	//don't filter it later on. Filter it in this function
	dontFilterPrefix = true

	path := parseFilepath(cur)
	searchPath := path

	filter := ""
	if !path.dir {
		searchPath.parts = searchPath.parts[:len(searchPath.parts)-1]
		searchPath.dir = true
		filter = path.parts[len(path.parts)-1]
	}

	log.Write("SEARCH PATH: %+v", searchPath.SearchString())
	contents, err := searchPath.GetContents(acceptFile)
	if err != nil {
		log.Write("Erred to get contents")
		return nil, err
	}

	if path.dir && len(contents) == 0 {
		log.Write("dir with no contents")
		dontAddSpace = false
		return []string{cur}, nil
	}

	log.Write("CONTENTS: %+v\n", contents)

	//Do our own filtering now
	candidates := []filepath{}

	if len(path.parts) > 0 {
		lastPart := path.parts[len(path.parts)-1]
		if lastPart == "." || lastPart == "./" || lastPart == ".." || lastPart == "../" {
			dotPath := make([]string, len(searchPath.parts))
			copy(dotPath, searchPath.parts)
			dotPath = append(dotPath, ".")
			candidates = append(candidates, filepath{parts: dotPath, dir: true, absolute: path.absolute})

			dotDotPath := make([]string, len(searchPath.parts))
			copy(dotDotPath, searchPath.parts)
			dotDotPath = append(dotDotPath, "..")
			candidates = append(candidates, filepath{parts: dotDotPath, dir: true, absolute: path.absolute})

			if lastPart == ".." || lastPart == "../" {
				candidates = candidates[1:]
			}
		}
	}

	for _, content := range contents {
		if !acceptFile && !content.dir {
			continue
		}

		if strings.HasPrefix(content.parts[len(content.parts)-1], filter) {
			candidates = append(candidates, content)
		}
	}

	//Check if we should kick out a space
	if len(candidates) == 1 {
		if !candidates[0].dir {
			dontAddSpace = false
		} else {
			nextContents, err := walkDirs(candidates[0].SearchString(), acceptFile)
			if err == nil && len(nextContents) == 0 { //Yes, should be == nil
				dontAddSpace = false
			}
		}
	}

	ret := []string{}
	for _, candidate := range candidates {
		ret = append(ret, candidate.OriginalString())
	}

	return ret, nil
}
