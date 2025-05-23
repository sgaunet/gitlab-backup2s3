package app_test

import (
	"testing"

	"github.com/sgaunet/gitlab-backup2s3/pkg/app"
)

func TestNewApp(t *testing.T) {
	a := app.NewApp()
	if a == nil {
		t.Fatalf("Expected non-nil app, got nil")
	}
	a.SetBackupCmd([]string{"echo", "hello"})
	err := a.Run()
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
}

func TestRunCommandReportsError(t *testing.T) {
	a := app.NewApp()
	if a == nil {
		t.Fatalf("Expected non-nil app, got nil")
	}
	a.SetBackupCmd([]string{"exit", "1"})
	err := a.Run()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestSetLoggerNil(t *testing.T) {
	a := app.NewApp()
	if a == nil {
		t.Fatalf("Expected non-nil app, got nil")
	}
	a.SetLogger(nil) // now accepts *logger.Logger as nil
	a.SetBackupCmd([]string{"echo", "hello"})
	err := a.Run()
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
}
