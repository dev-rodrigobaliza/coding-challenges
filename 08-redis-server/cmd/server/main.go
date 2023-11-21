package main

import (
	"flag"
	"fmt"
	"os"
	"rs/redis"
)

func main() {
	addr := flag.String("addr", ":6379", "server address to listen requests")
	log := flag.Bool("log", false, "enable server logging")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	redis := redis.New(*addr, *log)
	if err := redis.Start(); err != nil {
		panic(err)
	}
}
