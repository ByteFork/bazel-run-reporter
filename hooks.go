package main

import (
	"os"
	"os/exec"
	"strings"
)

type (
	command struct {
		raw   string
		parts []string
	}
)

func (c *command) String() string { return c.raw }
func (c *command) Set(s string) {
	c.parts = strings.Fields(s)
	c.raw = s
}

func (c *command) Get() []string {
	if c == nil {
		return nil
	}

	return c.parts
}

func (c *command) Execute() error {
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
