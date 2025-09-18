package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func Run(cmd *exec.Cmd) (string, error) {
	dir, _ := GetProjectDir()
	cmd.Dir = dir
	command := strings.Join(cmd.Args, " ")
	slog.Info("running shell command", "command", command, "dir", dir)

	if err := os.Chdir(cmd.Dir); err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("%q failed with error %q: %w", command, string(output), err)
	}

	return string(output), nil
}

func GetProjectDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return wd, fmt.Errorf("failed to get current working directory: %w", err)
	}
	wd = strings.ReplaceAll(wd, "/test", "")
	return wd, nil
}
