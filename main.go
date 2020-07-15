package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var authToken = os.Getenv("MINIO_WEBHOOK_AUTH_TOKEN")
var port = os.Getenv("MINIO_WEBHOOK_PORT")

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: minio-webhook filename.log")
		os.Exit(1)
	}
	if port == "" {
		port = "8080"
	}
	logFilePath := os.Args[1]
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatal(err)
	}
	var logFileMu sync.Mutex

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for _ = range sigs {
			logFileMu.Lock()
			logFile.Close()
			logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
			if err != nil {
				log.Fatal(err)
			}
			logFileMu.Unlock()
		}
	}()

	err = http.ListenAndServe(":"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authToken != "" {
			if authToken != r.Header.Get("Authorization") {
				return
			}
		}
		switch r.Method {
		case "POST":
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			logFileMu.Lock()
			logFile.Write(data)
			logFile.WriteString("\n")
			logFileMu.Unlock()
		default:
		}
	}))
	if err != nil {
		log.Fatal(err)
	}
}
