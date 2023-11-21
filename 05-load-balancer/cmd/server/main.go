package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	addr *string
)

func main() {
	addr = flag.String("addr", "", "server address to listen requests")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid server address")
		os.Exit(0)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/health", getHealth)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		panic(err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	msg := fmt.Sprintf("hello from backend server running at: %s\n", *addr)
	io.WriteString(w, msg)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /health request\n")

	io.WriteString(w, "ok")
}