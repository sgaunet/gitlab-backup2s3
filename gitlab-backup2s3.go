package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	var backupCmd []string
	debugLevel := os.Getenv("DEBUGLEVEL")
	initTrace(debugLevel)
	backupCmd = []string{
		"gitlab-backup",
	}
	log.Debugln("backupCmd=", backupCmd)
	log.Infoln("Execute gitlab-backup")
	err = execCommand(backupCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func execCommand(cmdToExec []string) error {
	cmd := exec.Command(cmdToExec[0], cmdToExec[1:]...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			log.Errorln(m)
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			log.Infoln(m)
		}
	}()
	err = cmd.Wait()
	return err
}

func initTrace(debugLevel string) {
	log.SetOutput(os.Stdout)
	switch debugLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
