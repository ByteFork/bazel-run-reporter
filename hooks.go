package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"

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

func (c *CommandHook) Execute(timeout time.Duration) error {
	command := c.Get()

	if len(command) == 0 {
		return nil
	}

	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// #nosec G204 - Command is from a trusted source (command line flag)
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err := cmd.Run()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return ctx.Err()
	}

	return err
}
