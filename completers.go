package main

import "fmt"

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

func compInstanceGroups(ctx compContext) ([]string, error) {
	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}
	instances, err := fetchInstances(client, ctx)
	if err != nil {
		return nil, err
	}
	uniqueMap := map[string]bool{}
	for _, instance := range instances {
		uniqueMap[instance.Job] = true
	}

	ret := make([]string, 0, len(uniqueMap))
	for group := range uniqueMap {
		ret = append(ret, group)
	}

	return ret, nil
}

func compInstances(ctx compContext) ([]string, error) {
	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}
	instances, err := fetchInstances(client, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(instances))
	for _, instance := range instances {
		ret = append(ret, fmt.Sprintf("%s/%s", instance.Job, instance.ID))
		ret = append(ret, fmt.Sprintf("%s/%d", instance.Job, instance.Index))
	}

	return ret, nil
}

func compReleases(ctx compContext) ([]string, error) {
	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}
	releases, err := fetchReleases(client)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, release := range releases {
		ret = append(ret, release.Name)
		for _, version := range release.Versions {
			ret = append(ret, fmt.Sprintf("%s/%s", release.Name, version.Version))
		}
	}

	return ret, nil
}

func compSpecificReleases(ctx compContext) ([]string, error) {
	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}
	releases, err := fetchReleases(client)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, release := range releases {
		for _, version := range release.Versions {
			ret = append(ret, fmt.Sprintf("%s/%s", release.Name, version.Version))
		}
	}

	return ret, nil
}

func compUnusedReleases(ctx compContext) ([]string, error) {
	client, err := getBoshClient(ctx)
	if err != nil {
		return nil, err
	}
	releases, err := fetchReleases(client)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, release := range releases {
		allUnused := true
		for _, version := range release.Versions {
			if version.CurrentlyDeployed {
				allUnused = false
			} else {
				ret = append(ret, fmt.Sprintf("%s/%s", release.Name, version.Version))
			}
		}

		if allUnused {
			ret = append(ret, release.Name)
		}
	}

	return ret, nil
}

func compOr(fns ...compFunc) compFunc {
	return func(ctx compContext) ([]string, error) {
		ret := []string{}
		for _, fn := range fns {
			theseComps, err := fn(ctx)
			if err != nil {
				return nil, err
			}

			ret = append(ret, theseComps...)
		}
		return ret, nil
	}
}
