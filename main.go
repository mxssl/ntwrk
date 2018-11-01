package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func main() {
	listenAddress := flag.String("listen", ":80", "what address and port to listen on (ex.: localhost:8080 or just :8080)")
	flag.Parse()
	log.Println("Starting the app...")
	http.HandleFunc("/", rootHandler)

	Server := http.Server{Addr: *listenAddress}

	go func() {
		log.Fatal(Server.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	switch <-signalChan {
	case syscall.SIGINT:
		log.Print("Gracefully shutting down")
		if err := Server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Cannot gracefully shut down the server: %v", err)
		}
	case syscall.SIGTERM:
		// abruptly shutting down if we got killed
		log.Println("Killed")
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// a bit lazy ipv6 regexp, but it should work well, since chances of getting unreal ipv6 values are 0
	re := regexp.MustCompile(`((?:\d{1,3}.){3}\d{1,3}):\d+|\[(.+)\]:\d+`)
	for _, ip := range re.FindStringSubmatch(r.RemoteAddr)[1:] {
		if ip != "" {
			fmt.Fprintf(w, ip+"\n")
			log.Printf("Received ip request from: %s\n", ip)
		}
	}
}
