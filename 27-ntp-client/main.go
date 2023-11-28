package main

import (
	"flag"
	"fmt"
	"ntpc/clock"
	"os"
)

func main() {
	ntp1 := flag.String("ntp1", "0.br.pool.ntp.org:123", "ntp server address to make queries")
	ntp2 := flag.String("ntp2", "1.br.pool.ntp.org:123", "ntp server address to make queries")

	flag.Parse()

	if *ntp1 == "" || *ntp2 == "" {
		fmt.Println("please specify two valid ntp server address")
		os.Exit(1)
	}

	t1, err := clock.QueryNTP(*ntp1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", t1)

	t2, err := clock.QueryNTP(*ntp2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", t2)
}
