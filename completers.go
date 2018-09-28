package main

func compNoop(ctx compContext) ([]string, error) {
	return nil, nil
}

func compCommandNames(ctx compContext) ([]string, error) {
	var ret []string
	for _, cmd := range commands {
		ret = append(ret, cmd.Name)
	}

	return ret, nil
}

func compFlagNames(ctx compContext) ([]string, error) {
	var ret []string
	for name := range flags {
		ret = append(ret, name)
	}

	return ret, nil
}

func compEnvAliases(ctx compContext) ([]string, error) {
	conf, err := getBoshConfig(ctx)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, env := range conf.Environments {
		ret = append(ret, env.Alias)
	}

	return ret, nil
}

func compEnum(s ...string) func(compContext) ([]string, error) {
	return func(compContext) ([]string, error) {
		return s, nil
	}
}

func compFiles(ctx compContext) ([]string, error) {
	return walkDirs(ctx.CurrentToken, true, false)
}

func compDirs(ctx compContext) ([]string, error) {
	return walkDirs(ctx.CurrentToken, false, true)
}

func compDeployments(ctx compContext) ([]string, error) {
	type deployment struct {
		Name string `json:"name"`
	}

	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}

	deployments := []deployment{}
	err = client.Get("/deployments", &deployments)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(deployments))
	for _, dep := range deployments {
		ret = append(ret, dep.Name)
	}

	return ret, nil
}
