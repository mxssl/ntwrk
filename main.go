package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	re   *regexp.Regexp
	mode string
	port string
)

const (
	modeNative     = "native"
	modeCloudflare = "cloudflare"
)

func init() {
	// Configure logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Load config from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Connot load .env file, env variables will be used")
	}

	// Assign env variables
	mode = os.Getenv("MODE")
	port = os.Getenv("PORT")

	// Compile regular expression
	re = regexp.MustCompile(`((?:\d{1,3}.){3}\d{1,3}):\d+|\[(.+)\]:\d+`)
}

func main() {
	flag.Parse()
	log.Println("Starting the app...")

	http.HandleFunc("/", rootHandler)

	port := fmt.Sprintf(":" + port)

	Server := http.Server{Addr: port}

	go func() {
		log.Fatal(Server.ListenAndServe())
	}()

	log.WithFields(log.Fields{
		"addr": Server.Addr,
	}).Info("App started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	switch <-signalChan {
	case syscall.SIGINT:
		log.Print("Gracefully shutting down")
		if err := Server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Cannot gracefully shut down the server: %v", err)
		}
	case syscall.SIGTERM:
		log.Println("Killed")
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if mode == modeNative {
		for _, ip := range re.FindStringSubmatch(r.RemoteAddr)[1:] {
			if ip != "" {
				fmt.Fprintf(w, ip+"\n")
				log.Printf("Received ip request from: %s", ip)
			}
		}
	}
	if mode == modeCloudflare {
		ip := r.Header.Get("CF-Connecting-IP")
		fmt.Fprintf(w, ip+"\n")
		log.Printf("Received ip request from: %s", ip)
	}
}
