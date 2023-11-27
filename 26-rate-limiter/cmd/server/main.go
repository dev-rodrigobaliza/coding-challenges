package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rl/api"
	"rl/limiter/algorithms"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8080", "server address to listen requests")
	dsn := flag.String("dsn", ":memory:", "store dsn")
	limit := flag.Int("limit", 10, "limit (requests per second)")
	log := flag.Bool("log", true, "log messages")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	tk := algorithms.NewTokenBucket(*dsn, *limit)
	api := api.New(tk, *log)
	api.Start(*addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	api.Stop()
}
