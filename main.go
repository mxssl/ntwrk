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

func init() {
	// Configure logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Load config from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Connot load .env file, env variables will be used")
	}
}

func main() {
	flag.Parse()
	log.Println("Starting the app...")
	http.HandleFunc("/", rootHandler)

	port := fmt.Sprintf(":" + os.Getenv("PORT"))

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
	if os.Getenv("MODE") == "native" {
		re := regexp.MustCompile(`((?:\d{1,3}.){3}\d{1,3}):\d+|\[(.+)\]:\d+`)
		for _, ip := range re.FindStringSubmatch(r.RemoteAddr)[1:] {
			if ip != "" {
				fmt.Fprintf(w, ip+"\n")
				log.Printf("Received ip request from: %s", ip)
			}
		}
	}
	if os.Getenv("MODE") == "cloudflare" {
		ip := r.Header.Get("CF-Connecting-IP")
		fmt.Fprintf(w, ip+"\n")
		log.Printf("Received ip request from: %s", ip)
	}
}
