package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/jcchavezs/pakay/internal/log"
)

// Executes shell a command
func Command(command string, args ...string) ([]byte, error) {
	return commandContext(context.Background(), os.Stderr, command, args...)
}

// Executes shell a command passing a context
func commandContext(ctx context.Context, stderr io.Writer, command string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = stderr

	log.Logger.Debug("Executing command", "command", cmd.String())

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s: %w", cmd.String(), err)
	}

	return out.Bytes(), nil
}

// Executes shell a command passing a context
func CommandContext(ctx context.Context, command string, args ...string) ([]byte, error) {
	return commandContext(ctx, os.Stderr, command, args...)
}

// Executes shell a command passing a context
func CommandContextQ(ctx context.Context, command string, args ...string) ([]byte, error) {
	return commandContext(ctx, io.Discard, command, args...)
}
