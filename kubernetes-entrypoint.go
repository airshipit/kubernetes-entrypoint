package main

import (
	"os"

	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/config"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/container"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/daemonset"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/job"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/pod"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/service"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/socket"
	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	command "opendev.org/airship/kubernetes-entrypoint/util/command"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

func main() {
	var comm []string
	var entrypoint *entry.Entrypoint
	var err error
	if entrypoint, err = entry.New(nil); err != nil {
		logger.Error.Printf("Creating entrypoint failed: %v", err)
		os.Exit(1)
	}

	entrypoint.Resolve()

	if comm = env.SplitCommand(); len(comm) == 0 {
		// TODO(DTadrzak): we should consider other options to handle whether pod
		// is an init-container
		logger.Warning.Printf("COMMAND env is empty")
		os.Exit(0)
	}

	if err = command.Execute(comm); err != nil {
		logger.Error.Printf("Cannot execute command: %v", err)
		os.Exit(1)
	}
}
