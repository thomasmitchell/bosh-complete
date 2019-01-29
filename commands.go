package main

import "sort"

type commandList []command

func (c commandList) Find(name string) (ret command, found bool) {
	left, right := 0, len(commands)
	log.Write("finding command: %s", name)

	for left < right {
		mid := (left + right) / 2
		if commands[mid].Name > name {
			right = mid
		} else if commands[mid].Name < name {
			left = mid + 1
		} else {
			ret = commands[mid]
			found = true
			break
		}
	}

	return
}

func (c *commandList) Populate() {
	command{
		Name: "add-blob",
	}.Insert()

	command{
		Name: "alias-env",
		Args: []compFunc{compEnvAliases},
	}.Insert()

	command{
		Name: "attach-disk",
		Flags: []flag{
			{Long: "disk-properties", Complete: compNoop},
		},
	}.Insert()

	command{
		Name:  "blobs",
		Flags: []flag{{Long: "dir", Complete: compDirs}},
	}.Insert()

	command{
		Name: "cancel-task",
	}.Insert().Alias("ct")

	command{
		Name:  "clean-up",
		Flags: []flag{{Long: "all"}},
	}.Insert()

	command{
		Name: "cloud-check",
		Flags: []flag{
			{Long: "auto", Short: 'a'},
			{Long: "report", Short: 'r'},
			{Long: "resolution", Complete: compNoop},
		},
	}.Insert().Alias("cck").Alias("cloudcheck")

	command{
		Name: "cloud-config",
	}.Insert().Alias("cc")

	command{
		Name: "config",
		Flags: []flag{
			//TODO: name -> config names
			{Long: "name", Complete: compNoop},
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi")},
		},
		Args: []compFunc{
			compNoop, //TODO: Config ids
		},
	}.Insert().Alias("c")

	command{
		Name: "configs",
		Flags: []flag{
			//TODO: name -> config names
			{Long: "name", Complete: compNoop},
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi")},
			{Long: "recent", Complete: compNoop},
		},
	}.Insert().Alias("cs")

	command{
		Name: "cpi-config",
	}.Insert()

	command{
		Name: "create-env",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			{Long: "skip-drain"},
			{Long: "state", Complete: compFiles},
			{Long: "recreate"},
			{Long: "recreate-persistent-disks"},
		},
		Args: []compFunc{compFiles},
	}.Insert()

	command{
		Name: "create-release",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
			{Long: "name", Complete: compNoop},
			{Long: "version", Complete: compNoop},
			{Long: "timestamp-version"},
			{Long: "final"},
			{Long: "tarball", Complete: compFiles},
			{Long: "force"},
		},
	}.Insert().Alias("cr")

	command{
		Name: "curl",
		Args: []compFunc{
			compNoop,
		},
	}.Insert()

	command{
		Name: "delete-config",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi")},
			//TODO: name -> config names
			{Long: "name", Complete: compNoop},
		},
		Args: []compFunc{
			compNoop, //TODO: Config ids
		},
	}.Insert().Alias("dc")

	command{
		Name:  "delete-deployment",
		Flags: []flag{{Long: "force"}},
	}.Insert().Alias("deld")

	command{
		Name: "delete-disk",
		Args: []compFunc{
			compNoop, //TODO: (orphaned?) Disk cids
		},
	}.Insert()

	command{
		Name: "delete-env",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			{Long: "skip-drain"},
			{Long: "state", Complete: compFiles},
		},
		Args: []compFunc{compFiles},
	}.Insert()

	command{
		Name: "delete-network",
		Args: []compFunc{
			compNoop, //TODO: Network names
		},
	}.Insert()

	command{
		Name:  "delete-release",
		Flags: []flag{{Long: "force"}},
		Args: []compFunc{
			compUnusedReleases,
		},
	}.Insert().Alias("delr")

	command{
		Name: "delete-snapshot",
		Args: []compFunc{
			compNoop, //TODO: snapshot cids
		},
	}.Insert()

	command{
		Name: "delete-snapshots",
	}.Insert().Alias("dels")

	command{
		Name:  "delete-stemcell",
		Flags: []flag{{Long: "force"}},
		Args: []compFunc{
			compNoop, //TODO: stemcell-name[/version]
		},
	}.Insert()

	command{
		Name: "delete-vm",
		Args: []compFunc{
			compNoop, //TODO: vm cids
		},
	}.Insert()

	command{
		Name: "deploy",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			{Long: "no-redact"},
			{Long: "recreate"},
			{Long: "recreate-persistent-disks"},
			{Long: "fix"},
			//TODO: skip-drain -> get instance groups from manifest
			{Long: "skip-drain", Complete: compNoop},
			{Long: "max-in-flight", Complete: compNoop},
			{Long: "dry-run"},
		},
		Args: []compFunc{compFiles},
	}.Insert().Alias("d")

	command{
		Name: "deployment",
	}.Insert().Alias("dep")

	command{
		Name: "deployments",
	}.Insert().Alias("ds")

	command{
		Name: "diff-config",
		Flags: []flag{
			//TODO: Config ids
			{Long: "from-id", Complete: compNoop},
			//TODO: Config ids
			{Long: "to-id", Complete: compNoop},
			{Long: "from-content", Complete: compFiles},
			{Long: "to-content", Complete: compFiles},
		},
	}.Insert()

	command{
		Name:  "disks",
		Flags: []flag{{Long: "orphaned"}},
	}.Insert()

	command{
		Name: "environment",
	}.Insert().Alias("env")

	command{
		Name: "environments",
	}.Insert().Alias("envs")

	command{
		Name: "errands",
	}.Insert().Alias("es")

	command{
		Name: "event",
		Args: []compFunc{
			compNoop, //TODO: Event IDs (?)
		},
	}.Insert()

	command{
		Name: "events",
		Flags: []flag{
			//TODO: Event IDs (?)
			{Long: "before-id", Complete: compNoop},
			{Long: "before", Complete: compNoop},
			{Long: "after", Complete: compNoop},
			//TODO: Task IDs (?)
			{Long: "task", Complete: compNoop},
			//TODO: Instances
			{Long: "instance", Complete: compNoop},
			//TODO: Event users (?)
			{Long: "event-user", Complete: compNoop},
			{Long: "action", Complete: compEnum("update", "delete", "setup ssh", "cleanup ssh")},
			{Long: "object-type", Complete: compEnum("instance", "deployment", "vm")},
			//TODO: Probably complete this, but only if object type is given?
			{Long: "object-name", Complete: compNoop},
		},
	}.Insert()

	command{
		Name: "export-release",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
			//TODO: List jobs in current release dir
			{Long: "job", Complete: compNoop},
		},
		Args: []compFunc{
			compNoop,
			compNoop,
		},
	}.Insert()

	command{
		Name: "finalize-release",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
			{Long: "name", Complete: compNoop},
			{Long: "version", Complete: compNoop},
			{Long: "force"},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert()

	command{
		Name:  "generate-job",
		Flags: []flag{{Long: "dir", Complete: compDirs}},
		Args:  []compFunc{compNoop},
	}.Insert()

	command{
		Name: "help",
	}.Insert()

	command{
		Name: "ignore",
		Args: []compFunc{
			compInstances,
		},
	}.Insert()

	command{
		Name: "init-release",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
			{Long: "git"},
		},
	}.Insert()

	command{
		Name: "inspect-local-stemcell",
		Args: []compFunc{compDirs},
	}.Insert()

	command{
		Name: "inspect-release",
		Args: []compFunc{compSpecificReleases},
	}.Insert()

	command{
		Name: "instances",
		Flags: []flag{
			{Long: "details", Short: 'i'},
			{Long: "dns"},
			{Long: "vitals"},
			{Long: "ps", Short: 'p'},
			{Long: "failing", Short: 'f'},
		},
	}.Insert().Alias("is")

	command{
		Name: "interpolate",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			//TODO: I think this is parsing the paths of a yaml file?
			{Long: "path", Complete: compNoop},
			{Long: "var-errs"},
			{Long: "var-errs-unused"},
		},
	}.Insert().Alias("int")

	command{
		Name: "locks",
	}.Insert()

	command{
		Name: "log-in",
	}.Insert().Alias("l").Alias("login")

	command{
		Name: "log-out",
	}.Insert().Alias("logout")

	command{
		Name: "logs",
		Flags: []flag{
			{Long: "dir", Complete: compNoop},
			{Long: "follow", Short: 'f'},
			{Long: "num", Complete: compNoop},
			{Long: "quiet", Short: 'q'},
			//TODO: Jobs on the VM
			{Long: "job", Complete: compNoop},
			{Long: "only", Complete: compNoop},
			{Long: "agent"},
			{Long: "gw-disable"},
			{Long: "gw-user", Complete: compNoop},
			{Long: "gw-host", Complete: compNoop},
			{Long: "gw-private-key", Complete: compFiles},
			{Long: "gw-socks5", Complete: compFiles},
		},
	}.Insert()

	command{
		Name: "manifest",
	}.Insert().Alias("man")

	command{
		Name:  "networks",
		Flags: []flag{{Long: "orphaned", Short: 'o'}},
	}.Insert()

	command{
		Name: "orphan-disk",
		Args: []compFunc{
			compNoop, //TODO: disk cids (non-orphaned)
		},
	}.Insert()

	command{
		Name: "orphaned-vms",
	}.Insert()

	command{
		Name: "recreate",
		Flags: []flag{
			{Long: "skip-drain"},
			{Long: "force"},
			{Long: "fix"},
			{Long: "canaries", Complete: compNoop},
			{Long: "max-in-flight", Complete: compNoop},
			{Long: "dry-run"},
		},
		Args: []compFunc{
			compOr(compInstanceGroups, compInstances),
		},
	}.Insert()

	command{
		Name: "releases",
	}.Insert().Alias("rs")

	command{
		Name:  "remove-blob",
		Flags: []flag{{Long: "dir", Complete: compDirs}},
		Args: []compFunc{
			compNoop, //TODO: Not sure if file path or path within blob registry (i.e. blob name)
		},
	}.Insert()

	command{
		Name: "repack-stemcell",
		Flags: []flag{
			{Long: "name", Complete: compNoop},
			{Long: "cloud-properties", Complete: compNoop},
			{Long: "empty-image"},
			{Long: "format", Complete: compNoop},
			{Long: "version", Complete: compNoop},
		},
	}.Insert()

	command{
		Name:  "reset-release",
		Flags: []flag{{Long: "dir", Complete: compDirs}},
	}.Insert()

	command{
		Name: "restart",
		Flags: []flag{
			{Long: "skip-drain"},
			{Long: "force"},
			{Long: "canaries", Complete: compNoop},
			{Long: "max-in-flight", Complete: compNoop},
		},
	}.Insert()

	command{
		Name: "run-errand",
		Flags: []flag{
			{Long: "instance", Complete: compInstances},
			{Long: "keep-alive"},
			{Long: "when-changed"},
			{Long: "download-logs"},
			{Long: "logs-dir", Complete: compDirs},
		},
	}.Insert()

	command{
		Name: "runtime-config",
		Flags: []flag{
			//TODO: Probably the name of runtime configs?
			{Long: "name", Complete: compNoop},
		},
	}.Insert().Alias("rc")

	command{
		Name: "scp",
		Flags: []flag{
			{Long: "recursive", Short: 'r'},
			{Long: "gw-disable"},
			{Long: "gw-user", Complete: compNoop},
			{Long: "gw-host", Complete: compNoop},
			{Long: "gw-private-key", Complete: compFiles},
			{Long: "gw-socks5", Complete: compFiles},
		},
		Args: []compFunc{
			compFiles, //TODO: at least instance group/id... maybe use ssh to ls if thats not too slow?
			//TODO: "or" that with files on the local file system
			compFiles,
		},
	}.Insert()

	command{
		Name: "snapshots",
		Args: []compFunc{
			compInstances,
		},
	}.Insert()

	command{
		Name: "ssh",
		Flags: []flag{
			{Long: "command", Short: 'c', Complete: compNoop},
			{Long: "opts", Complete: compNoop},
			{Long: "results", Short: 'r'},
			{Long: "gw-disable"},
			{Long: "gw-user", Complete: compNoop},
			{Long: "gw-host", Complete: compNoop},
			{Long: "gw-private-key", Complete: compFiles},
			{Long: "gw-socks5", Complete: compFiles},
		},
		Args: []compFunc{
			compInstances,
		},
	}.Insert()

	command{
		Name: "start",
		Flags: []flag{
			{Long: "force"},
			{Long: "canaries", Complete: compNoop},
			{Long: "max-in-flight", Complete: compNoop},
		},
		Args: []compFunc{
			compOr(compInstanceGroups, compInstances),
		},
	}.Insert()

	command{
		Name: "stemcells",
	}.Insert().Alias("ss")

	command{
		Name: "stop",
		Flags: []flag{
			{Long: "soft"},
			{Long: "hard"},
			{Long: "skip-drain"},
			{Long: "force"},
			{Long: "canaries", Complete: compNoop},
			{Long: "max-in-flight", Complete: compNoop},
		},
		Args: []compFunc{
			compOr(compInstanceGroups, compInstances),
		},
	}.Insert()

	command{
		Name:  "sync-blobs",
		Flags: []flag{{Long: "dir", Complete: compDirs}},
	}.Insert()

	command{
		Name: "take-snapshot",
		Args: []compFunc{
			compInstances,
		},
	}.Insert()

	command{
		Name: "task",
		Flags: []flag{
			{Long: "event"},
			{Long: "cpi"},
			{Long: "debug"},
			{Long: "result"},
			{Long: "all", Short: 'a'},
		},
		Args: []compFunc{
			compNoop, //TODO: Task ids?
		},
	}.Insert().Alias("t")

	command{
		Name: "tasks",
		Flags: []flag{
			{Long: "recent", Complete: compNoop},
			{Long: "all", Short: 'a'},
		},
	}.Insert().Alias("ts")

	command{
		Name: "unignore",
		Args: []compFunc{
			compInstances,
		},
	}.Insert()

	command{
		Name: "update-cloud-config",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert().Alias("ucc")

	command{
		Name: "update-config",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi")},
			//TODO: Config names with type --type
			{Long: "name", Complete: compNoop},
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert().Alias("uc")

	command{
		Name: "update-cpi-config",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			{Long: "no-redact"},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert()

	command{
		Name: "update-resurrection",
		Args: []compFunc{
			compEnum("on", "off"),
		},
	}.Insert()

	command{
		Name: "update-runtime-config",
		Flags: []flag{
			//TODO: var -> <vars in manifest> = noop
			{Long: "var", Short: 'v', Complete: compNoop},
			//TODO: var-file -> <vars in manifest> = path
			{Long: "var-file", Complete: compNoop},
			{Long: "vars-file", Short: 'l', Complete: compFiles},
			{Long: "vars-env", Complete: compNoop},
			{Long: "vars-store", Complete: compFiles},
			{Long: "ops-file", Short: 'o', Complete: compFiles},
			{Long: "no-redact"},
			//TODO: Runtime config names
			{Long: "name", Complete: compNoop},
		},
	}.Insert().Alias("urc")

	command{
		Name: "upload-blobs",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
		},
	}.Insert()

	command{
		Name: "upload-release",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
			{Long: "rebase"},
			{Long: "fix"},
			{Long: "name", Complete: compNoop},
			{Long: "version", Complete: compNoop},
			{Long: "sha1", Complete: compNoop},
			{Long: "stemcell", Complete: compNoop},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert().Alias("ur")

	command{
		Name: "upload-stemcell",
		Flags: []flag{
			{Long: "fix"},
			{Long: "name", Complete: compNoop},
			{Long: "version", Complete: compNoop},
			{Long: "sha1", Complete: compNoop},
		},
		Args: []compFunc{
			compFiles,
		},
	}.Insert().Alias("us")

	command{
		Name: "variables",
	}.Insert().Alias("vars")

	command{
		Name: "vendor-package",
		Flags: []flag{
			{Long: "dir", Complete: compDirs},
		},
		Args: []compFunc{
			compNoop, //if the args were reversed, I could search the release dir for packages
			compDirs,
		},
	}.Insert()

	command{
		Name: "vms",
		Flags: []flag{
			{Long: "dns"},
			{Long: "vitals"},
			{Long: "cloud-properties"},
		},
	}.Insert()

	sort.Slice(commands, func(i, j int) bool { return commands[i].Name < commands[j].Name })
}

var commands commandList

type command struct {
	Name  string
	Flags []flag
	Args  []compFunc
}

func (c command) Insert() command {
	commands = append(commands, c)
	return c
}

func (c command) Alias(alias string) command {
	c.Name = alias
	return c.Insert()
}

func (c command) InsertFlags() {
	for _, f := range c.Flags {
		insertFlag(f)
	}
}
