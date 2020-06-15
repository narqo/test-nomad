package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	version   string
	commit    string
	buildTime string
)

const serverShutdownTimeout = 5 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigc
		signal.Stop(sigc)
		cancel()
	}()

	if err := run(ctx, os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	printVersion := flags.Bool("version", false, "print version and exit")

	var (
		addr         string
		redisConnURL string
		redisDB      string
	)
	flags.StringVar(&addr, "http.addr", "127.0.0.1:8080", "address to listen on")
	flags.StringVar(&redisConnURL, "redis.conn-url", "tcp://127.0.0.1:6379", "redis connection url")
	flags.StringVar(&redisDB, "redis.db", "127.0.0.1:6379", "redis connection address")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *printVersion {
		log.Printf("version %s", versionString())
		os.Exit(1)
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/internal/ping", pingHandler)
	http.HandleFunc("/internal/stats", statsHandler)
	http.HandleFunc("/internal/die", dieHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	time.Sleep(2 * time.Second) // simulate slow start

	errc := make(chan error, 1)
	go func() {
		log.Printf("starting: addr %s, version %s", addr, versionString())
		errc <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("exiting...")
	case err := <-errc:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	return server.Shutdown(ctx)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello!")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "pong")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := struct {
		Tag      string `json:"tag"`
		NodeID   string `json:"node_id"`
		NodeName string `json:"node_name"`
		Env      string `json:"env"`
		Version  string `json:"version"`
	}{
		Tag:      os.Getenv("ADJUST_TAG"),
		NodeID:   os.Getenv("ADJUST_NODE_ID"),
		NodeName: os.Getenv("ADJUST_NODE_NAME"),
		Env:      os.Getenv("GO_ENV"),
		Version:  versionString(),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("stats handler: failed encoding response: %v", err)
	}
}

// simulate unhandled panic
func dieHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		log.Panic("die handler: (un)expected panic")
	}()
}

func versionString() string {
	return fmt.Sprintf("%s, commit %s (%s), go version %s",
		version,
		commit,
		buildTime,
		runtime.Version(),
	)
}
