package main

import (
	"flag"
	"fmt"
	"os"
	"ws/ws"
)

func main() {
	addr := flag.String("addr", ":80", "server address to listen requests")
	dir := flag.String("dir", "", "directory to store html files")
	log := flag.Bool("log", false, "enable server logging")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(1)
	}

	webServer := ws.New(*addr, *dir, *log)
	if err := webServer.Start(); err != nil {
		panic(err)
	}
}
