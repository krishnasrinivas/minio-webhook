package main

import (
	"fmt"

	"github.com/kjk/dailyrotate"
)

func onLogClose(path string, didRotate bool) {
	fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
	if !didRotate {
		return
	}
	// process just closed file e.g. upload to backblaze storage for backup
	go func() {
		// if processing takes a long time, do it in background
	}()
}

var (
	logFile *dailyrotate.File
)

func openLogFile(pathFormat string, onClose func(string, bool)) error {
	w, err := dailyrotate.NewFile(pathFormat, onLogClose)
	if err != nil {
		return err
	}
	logFile = w
	return nil
}

func closeLogFile() error {
	return logFile.Close()
}

func writeToLog(buf []byte) error {
	_, err := logFile.Write(buf)
	return err
}
