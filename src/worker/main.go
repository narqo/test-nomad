package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

	var workerPollInterval time.Duration
	flags.DurationVar(&workerPollInterval, "worker.poll-internal", 3*time.Second, "worker poll interval")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *printVersion {
		log.Printf("version %s", versionString())
		os.Exit(1)
	}

	log.Printf("starting: version %s", versionString())

	return runWorker(ctx, workerPollInterval)
}

func runWorker(ctx context.Context, pollInterval time.Duration) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p := payloadBytes()
			log.Printf("working, %s\n", p)
		case <-ctx.Done():
			return nil
		}
	}
}

func payloadBytes() []byte {
	stats := struct {
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
	p, _ := json.Marshal(stats)
	return p
}

func versionString() string {
	return fmt.Sprintf("%s, commit %s (%s), go version %s",
		version,
		commit,
		buildTime,
		runtime.Version(),
	)
}
