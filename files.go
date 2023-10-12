package main

import "os"

func isDirExists(dir string) bool {
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	i, _ := os.Stat(dir)
	return i.IsDir()
}

func isFileExists(file string) bool {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	i, _ := os.Stat(file)
	return !i.IsDir()
}

func isThereAFileOrDir(file string) bool {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return false
	}
	defer f.Close()
	return true
}
