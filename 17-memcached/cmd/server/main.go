package main

import (
	"errors"
	"flag"
	"fmt"
	"mc/memcached"
	"mc/safemap"
	"mc/store"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":11211", "server address to listen requests")
	dsn := flag.String("dsn", ":memory:", "store dsn")
	log := flag.Bool("log", true, "log messages")

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

	mc := memcached.New(*addr, data, *log)
	if err := mc.Start(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := mc.Stop(); err != nil {
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
