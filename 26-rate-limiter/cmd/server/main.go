package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rl/api"
	"rl/limiter/algorithms"
	"rl/safemap"
	"rl/store"
	"syscall"
	"time"
)

func main() {
	addr := flag.String("addr", ":8080", "server address to listen requests")
	dsn := flag.String("dsn", ":memory:", "store dsn")
	limit := flag.Int("limit", 10, "max requests")
	dur := flag.String("dur", "1s", "limiter duration")
	log := flag.Bool("log", true, "log messages")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	// tk := getTokenBucket(*dsn, *dur, *limit)
	fw := getFixedWindow(*dsn, *dur, *limit)
	api := api.New(fw, *log)
	api.Start(*addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	api.Stop()
}

func getDuration(dur string) time.Duration {
	d, err := time.ParseDuration(dur)
	if err != nil {
		fmt.Println("please specify a valid limiter duration")
		os.Exit(1)
	}

	return d
}

func getStore(dsn string) store.Store {
	switch dsn {
	case ":memory:":
		return safemap.New()

	default:
		fmt.Println("please specify a valid dsn for the store")
		os.Exit(1)
	}

	return nil
}

func getTokenBucket(dsn, dur string, limit int) *algorithms.TokenBucket {
	s := getStore(dsn)
	d := getDuration(dur)
	return algorithms.NewTokenBucket(s, limit, d)
}

func getFixedWindow(dsn, dur string, limit int) *algorithms.FixedWindow {
	s := getStore(dsn)
	d := getDuration(dur)
	return algorithms.NewFixedWindow(s, limit, d)
}