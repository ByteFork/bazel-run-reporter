package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestCommandSet(t *testing.T) {
	tests := []struct {
		name      string
		cmdStr    string
		wantParts []string
	}{
		{
			name:      "simple command",
			cmdStr:    "echo hello",
			wantParts: []string{"echo", "hello"},
		},
		{
			name:      "command with quotes",
			cmdStr:    "echo \"hello world\"",
			wantParts: []string{"echo", "hello world"},
		},
		{
			name:      "command with single quotes",
			cmdStr:    "echo 'hello world'",
			wantParts: []string{"echo", "hello world"},
		},
		{
			name:      "command with mixed quotes",
			cmdStr:    "curl -X POST -F \"file=@results.xml\" http://localhost:8080",
			wantParts: []string{"curl", "-X", "POST", "-F", "file=@results.xml", "http://localhost:8080"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandHook{}
			if err := c.Set(tt.cmdStr); err != nil {
				t.Errorf("command.Set() error = %v", err)
				return
			}

			if c.String() != tt.cmdStr {
				t.Errorf("command.String() = %v, want %v", c.String(), tt.cmdStr)
			}

			parts := c.Get()
			if len(parts) != len(tt.wantParts) {
				t.Errorf("command.Get() length = %v, want %v", len(parts), len(tt.wantParts))
				return
			}

			for i, part := range parts {
				if part != tt.wantParts[i] {
					t.Errorf("command.Get()[%d] = %v, want %v", i, part, tt.wantParts[i])
				}
			}
		})
	}
}

func TestCommandExecute(t *testing.T) {
	if os.Getenv("TEST_COMMAND_EXECUTE") != "1" {
		t.Skip("Skipping execution test. Set TEST_COMMAND_EXECUTE=1 to run")
	}

	t.Run("echo command", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		c := &CommandHook{}
		if err := c.Set("echo test execution"); err != nil {
			t.Errorf("command.Set() error = %v", err)
			return
		}

		err := c.Execute()

		w.Close()

		os.Stdout = oldStdout

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Errorf("command.Execute() error = %v", err)
		}

		if output != "test execution\n" {
			t.Errorf("command.Execute() output = %q, want %q", output, "test execution\n")
		}
	})

	t.Run("failing command", func(t *testing.T) {
		c := &CommandHook{}
		if err := c.Set("command_that_does_not_exist"); err != nil {
			t.Errorf("command.Set() error = %v", err)
			return
		}

		err := c.Execute()

		if err == nil {
			t.Error("command.Execute() expected error for non-existent command")
		}

		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Errorf("command.Execute() error type = %T, want %T", err, &exec.ExitError{})
		}
	})
}

func TestEmptyCommand(t *testing.T) {
	var (
		nilCmd *CommandHook
		c      = &CommandHook{}
	)

	if parts := nilCmd.Get(); parts != nil {
		t.Errorf("nil command.Get() = %v, want nil", parts)
	}

	if err := c.Set(""); err != nil {
		t.Errorf("command.Set() error = %v", err)
		return
	}

	if err := c.Execute(); err != nil {
		t.Errorf("empty command.Execute() error = %v, want nil", err)
	}
}
