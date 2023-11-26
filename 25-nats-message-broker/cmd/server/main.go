package main

import (
	"flag"
	"fmt"
	"nats/engine"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":4222", "server address to listen requests")
	dsn := flag.String("dsn", ":memory:", "store dsn")
	log := flag.Bool("log", true, "enable server logging")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	nat := engine.New(*addr, *dsn, *log)
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
