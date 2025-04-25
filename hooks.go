package main

import (
	"os"
	"os/exec"

	"github.com/google/shlex"
)

type (
	CommandHook struct {
		cmd   string
		parts []string
	}
)

func (c *CommandHook) String() string { return c.cmd }
func (c *CommandHook) Set(s string) error {
	var err error

	c.parts, err = shlex.Split(s)
	c.cmd = s

	return err
}

func (c *CommandHook) Get() []string {
	if c == nil {
		return nil
	}

	return c.parts
}

func (c *CommandHook) Execute() error {
	command := c.Get()

	if len(command) == 0 {
		return nil
	}

	// #nosec G204 - Command is from a trusted source (command line flag)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}
