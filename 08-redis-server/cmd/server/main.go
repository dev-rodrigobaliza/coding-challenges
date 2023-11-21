package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"rs/redis"
	"rs/safemap"
	"rs/store"
)

func main() {
	addr := flag.String("addr", ":6379", "server address to listen requests")
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

	redis := redis.New(*addr, data, *log)
	if err := redis.Start(); err != nil {
		panic(err)
	}
}

func getStore(dsn string) (store.Store, error) {
	if dsn == ":memory:" {
		data := safemap.New()
		return data, nil
	}

	return nil, errors.New("unknown store")
}
