// Package app contains the main application logic for gitlab-backup2s3.
package app

import (
	"bufio"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/sgaunet/gitlab-backup2s3/pkg/logger"
)

// App is the main application.
type App struct {
	logger    logger.Logger
	backupCmd []string
}

// NewApp creates a new App.
func NewApp() *App {
	return &App{
		logger: logger.NoLogger(),
		backupCmd: []string{
			"gitlab-backup",
		},
	}
}

// SetLogger sets the logger.
func (a *App) SetLogger(log logger.Logger) {
	if log == nil {
		a.logger = logger.NoLogger()
		return
	}
	a.logger = log
}

// SetBackupCmd sets the backup command.
// Use this method for testing purposes.
func (a *App) SetBackupCmd(backupCmd []string) {
	a.backupCmd = backupCmd
}

// Run executes the application.
func (a *App) Run() error {
	a.logger.Info("Execute gitlab-backup")
	err := a.execCommand(a.backupCmd)
	if err != nil {
		a.logger.Error("Error executing gitlab-backup", slog.String("error", err.Error()))
		return err
	}
	return nil
}

// execCommand executes a command.
// It wraps all errors from external packages.
func (a *App) execCommand(cmdToExec []string) error {
	cmd := exec.Command(cmdToExec[0], cmdToExec[1:]...) // #nosec G204
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			a.logger.Error(m)
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			a.logger.Info(m)
		}
	}()
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for command: %w", err)
	}
	return nil
}
