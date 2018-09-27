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
	conf, err := getBoshConfig()
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
