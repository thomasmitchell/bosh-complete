package main

import "sort"

type commandList []command

func (c commandList) Find(name string) (ret command, found bool) {
	left, right := 0, len(commands)
	log.Write("finding command: %s", name)

	for left < right {
		mid := (left + right) / 2
		log.Write("left: %d, right: %d, mid: %d", left, right, mid)
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
	}.Insert()

	command{
		Name: "attach-disk",
	}.Insert()

	command{
		Name: "blobs",
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
		},
	}.Insert().Alias("cck").Alias("cloudcheck")

	command{
		Name: "cloud-config",
	}.Insert().Alias("cc")

	command{
		Name: "config",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi"), TakesValue: true},
		},
	}.Insert().Alias("c")

	command{
		Name: "configs",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi"), TakesValue: true},
		},
	}.Insert().Alias("cs")

	command{
		Name: "cpi-config",
	}.Insert()

	command{
		Name: "create-env",
	}.Insert()

	command{
		Name: "create-release",
		Flags: []flag{
			{Long: "timestamp-version"},
			{Long: "final"},
			{Long: "force"},
		},
	}.Insert().Alias("cr")

	command{
		Name: "delete-config",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi"), TakesValue: true},
		},
	}.Insert().Alias("dc")

	command{
		Name: "delete-deployment",
		Flags: []flag{
			{Long: "force"},
		},
	}.Insert().Alias("deld")

	command{
		Name: "delete-disk",
	}.Insert()

	command{
		Name: "delete-env",
		Flags: []flag{
			{Long: "skip-drain"},
		},
	}.Insert()

	command{
		Name: "delete-network",
	}.Insert()

	command{
		Name:  "delete-release",
		Flags: []flag{{Long: "force"}},
	}.Insert().Alias("delr")

	command{
		Name: "delete-snapshot",
	}.Insert()

	command{
		Name: "delete-snapshots",
	}.Insert().Alias("dels")

	command{
		Name:  "delete-stemcell",
		Flags: []flag{{Long: "force"}},
	}.Insert()

	command{
		Name: "delete-vm",
	}.Insert()

	command{
		Name: "deploy",
		Flags: []flag{
			{Long: "no-redact"},
			{Long: "recreate"},
			{Long: "recreate-persistent-disks"},
			{Long: "fix"},
			{Long: "dry-run"},
		},
	}.Insert().Alias("d")

	command{
		Name: "deployment",
	}.Insert().Alias("dep")

	command{
		Name: "deployments",
	}.Insert().Alias("ds")

	command{
		Name: "diff-config",
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
	}.Insert()

	command{
		Name: "events",
		Flags: []flag{
			{Long: "action", Complete: compEnum("update", "delete", "setup ssh", "cleanup ssh"), TakesValue: true},
		},
	}.Insert()

	command{
		Name: "export-release",
	}.Insert()

	command{
		Name: "finalize-release",
		Flags: []flag{
			{Long: "force"},
		},
	}.Insert()

	command{
		Name: "generate-job",
	}.Insert()

	command{
		Name: "help",
	}.Insert()

	command{
		Name: "ignore",
	}.Insert()

	command{
		Name: "init-release",
		Flags: []flag{
			{Long: "git"},
		},
	}.Insert()

	command{
		Name: "inspect-local-stemcell",
	}.Insert()

	command{
		Name: "inspect-release",
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
			{Long: "follow", Short: 'f'},
			{Long: "quiet", Short: 'q'},
			{Long: "agent"},
			{Long: "gw-disable"},
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
			{Long: "dry-run"},
		},
	}.Insert()

	command{
		Name: "releases",
	}.Insert().Alias("rs")

	command{
		Name: "remove-blob",
	}.Insert()

	command{
		Name: "repack-stemcell",
		Flags: []flag{
			{Long: "empty-image"},
		},
	}.Insert()

	command{
		Name: "reset-release",
	}.Insert()

	command{
		Name: "restart",
		Flags: []flag{
			{Long: "skip-drain"},
			{Long: "force"},
		},
	}.Insert()

	command{
		Name: "run-errand",
		Flags: []flag{
			{Long: "keep-alive"},
			{Long: "when-changed"},
			{Long: "download-logs"},
		},
	}.Insert()

	command{
		Name: "runtime-config",
	}.Insert().Alias("rc")

	command{
		Name: "scp",
		Flags: []flag{
			{Long: "recursive", Short: 'r'},
			{Long: "gw-disable"},
		},
	}.Insert()

	command{
		Name: "snapshots",
	}.Insert()

	command{
		Name: "ssh",
		Flags: []flag{
			{Long: "results", Short: 'r'},
			{Long: "gw-disable"},
		},
	}.Insert()

	command{
		Name: "start",
		Flags: []flag{
			{Long: "force"},
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
		},
	}.Insert()

	command{
		Name: "sync-blobs",
	}.Insert()

	command{
		Name: "take-snapshot",
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
	}.Insert().Alias("t")

	command{
		Name: "tasks",
		Flags: []flag{
			{Long: "all", Short: 'a'},
		},
	}.Insert().Alias("ts")

	command{
		Name: "unignore",
	}.Insert()

	command{
		Name: "update-cloud-config",
	}.Insert().Alias("ucc")

	command{
		Name: "update-config",
		Flags: []flag{
			{Long: "type", Complete: compEnum("cloud", "runtime", "cpi"), TakesValue: true},
		},
	}.Insert().Alias("uc")

	command{
		Name: "update-cpi-config",
		Flags: []flag{
			{Long: "no-redact"},
		},
	}.Insert()

	command{
		Name: "update-resurrection",
		ArgComps: []compFunc{
			compEnum("on", "off"),
		},
	}.Insert()

	command{
		Name: "update-runtime-config",
		Flags: []flag{
			{Long: "no-redact"},
		},
	}.Insert().Alias("urc")

	command{
		Name: "upload-blobs",
	}.Insert()

	command{
		Name: "upload-release",
		Flags: []flag{
			{Long: "rebase"},
			{Long: "fix"},
		},
	}.Insert().Alias("ur")

	command{
		Name: "upload-stemcell",
		Flags: []flag{
			{Long: "fix"},
		},
	}.Insert().Alias("us")

	command{
		Name: "variables",
	}.Insert().Alias("vars")

	command{
		Name: "vendor-package",
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
	Name     string
	Flags    []flag
	ArgComps []compFunc
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
