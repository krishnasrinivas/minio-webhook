package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var authToken = os.Getenv("MINIO_WEBHOOK_AUTH_TOKEN")
var port = os.Getenv("MINIO_WEBHOOK_PORT")

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: minio-webhook log-dir")
		os.Exit(1)
	}
	if port == "" {
		port = "8080"
	}
	logDir := os.Args[1]

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("os.MkdirAll() %s\n", err)
	}

	pathFormat := filepath.Join(logDir, "2006-01-02.log")
	if err := openLogFile(pathFormat, onLogClose); err != nil {
		log.Fatalf("openLogFile failed with '%s'\n", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for range sigs {
			closeLogFile()
		}
	}()

	err := http.ListenAndServe(":"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println(string(data))
			writeToLog(data)
		default:
		}
	}))
	if err != nil {
		log.Fatal(err)
	}
}
