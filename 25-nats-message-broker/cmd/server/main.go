package main

import (
	"errors"
	"flag"
	"fmt"
	"nats/engine"
	"nats/safemap"
	"nats/store"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":4222", "server address to listen requests")
	dsn := flag.String("dsn", ":memory:", "store dsn")
	log := flag.Bool("log", false, "enable server logging")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	data, err := getStore(*dsn)
	if err != nil {
		fmt.Println("failed to create the store", err)
		os.Exit(1)
	}

	nat := engine.New(*addr, data, *log)
	if err := nat.Start(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := nat.Stop(); err != nil {
		os.Exit(1)
	}
}

func getStore(dsn string) (store.Store, error) {
	if dsn == ":memory:" {
		data := safemap.New()
		return data, nil
	}

	return nil, errors.New("unknown store")
}
