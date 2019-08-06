package command

import (
	"os"
	"os/exec"
	"syscall"

	"opendev.org/airship/kubernetes-entrypoint/logger"
)

func Execute(command []string) error {
	path, err := exec.LookPath(command[0])
	if err != nil {
		logger.Error.Printf("Cannot find a binary %v : %v", command[0], err)
		return err
	}

	env := os.Environ()
	err = syscall.Exec(path, command, env)
	if err != nil {
		logger.Error.Printf("Executing command %v failed: %v", command, err)
		return err
	}
	return nil
}
