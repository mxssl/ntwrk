package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"log/slog"
)

var (
	re   *regexp.Regexp
	mode string
	port string
)

const (
	modeNative     = "native"
	modeCloudflare = "cloudflare"
	modeProxy      = "proxy"
)

func main() {
	// Configure logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load config from .env
	if err := godotenv.Load(); err != nil {
		slog.Info("cannot load .env file, env variables will be used")
	}

	// Assign env variables
	mode = os.Getenv("MODE")
	port = os.Getenv("PORT")

	// Compile regular expression
	re = regexp.MustCompile(`((?:\d{1,3}.){3}\d{1,3}):\d+|\[(.+)\]:\d+`)

	// Parse flags
	flag.Parse()

	// Start the app
	slog.Info("starting the app...")

	http.HandleFunc("/", rootHandler)

	port := fmt.Sprintf(":%s", port)

	server := http.Server{Addr: port}

	slog.Info(
		"app started",
		slog.String("addr", server.Addr),
	)
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "errorMsg", err)
		}
		return
	case sig := <-signalChan:
		slog.Info("shutting down...", "signal", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"cannot gracefully shutdown the app",
			"errorMsg", err,
		)
		if err := server.Close(); err != nil {
			slog.Error("cannot force the server to close", "errorMsg", err)
		}
	}

	if err := <-serverErrors; err != nil && err != http.ErrServerClosed {
		slog.Error("server error during shutdown", "errorMsg", err)
	}

	slog.Info("app stopped")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if mode == modeNative {
		for _, ip := range re.FindStringSubmatch(r.RemoteAddr)[1:] {
			if ip != "" {
				if _, err := fmt.Fprintf(w, "%s\n", ip); err != nil {
					slog.Error("failed to write response", "error", err)
					return
				}
				slog.Info(
					"request",
					slog.String("source", ip),
				)
			}
		}
	}
	if mode == modeCloudflare {
		ip := r.Header.Get("CF-Connecting-IP")
		if _, err := fmt.Fprintf(w, "%s\n", ip); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
		slog.Info(
			"request",
			"ip", ip,
		)
	}
	if mode == modeProxy {
		xff := r.Header.Get("X-Forwarded-For")
		// Extract first IP from comma-separated list
		ip := strings.Split(xff, ",")[0]
		ip = strings.TrimSpace(ip)

		if _, err := fmt.Fprintf(w, "%s\n", ip); err != nil {
			slog.Error("failed to write response", "error", err)
			return
		}
		slog.Info(
			"request",
			"ip", ip,
		)
	}
}
