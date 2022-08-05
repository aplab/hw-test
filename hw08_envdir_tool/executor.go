package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var (
	ErrUnableToSetVar   = errors.New("unable to set var")
	ErrUnableToUnsetVar = errors.New("unable to unset var")
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	err := setVariables(env)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
		return 1
	}
	command.Env = os.Environ()
	if err := command.Run(); err != nil {
		var ee *exec.ExitError
		_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
		if errors.As(err, &ee) {
			return ee.ExitCode()
		}
		return 1
	}
	return 0
}

func setVariables(env Environment) error {
	for name, value := range env {
		err := setVariable(name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func setVariable(name string, value EnvValue) error {
	if value.NeedRemove {
		err := os.Unsetenv(name)
		if err != nil {
			return ErrUnableToUnsetVar
		}
		return nil
	}
	err := os.Setenv(name, value.Value)
	if err != nil {
		return ErrUnableToSetVar
	}
	return nil
}
