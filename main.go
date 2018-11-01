package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	log.Println("Starting the app...")

	http.HandleFunc("/", rootHandler)

	Server := http.Server{Addr: ":80"}
	go func() {
		log.Fatal(Server.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
	Server.Shutdown(context.Background())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	addr := strings.Split(r.RemoteAddr, ":")
	log.Printf("Received ip request from: %v\n", addr[0])
	fmt.Fprintf(w, addr[0]+"\n")
}
