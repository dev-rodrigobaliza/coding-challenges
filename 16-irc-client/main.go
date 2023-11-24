package main

import (
	"flag"
	"fmt"
	"mirc/manager"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", "irc.freenode.net", "irc server address to connect")
	port := flag.String("port", "6697", "irc server port to connect")
	ssl := flag.Bool("ssl", true, "use a secure connection")
	username := flag.String("username", "gopher", "the username")
	realname := flag.String("realname", "gopher", "the realname")
	nickname := flag.String("nickname", "pupilo_do_sr_reitor", "the nickname")
	password := flag.String("password", "", "the password (optional)")
	log := flag.Bool("log", true, "log messages")

	flag.Parse()

	if *addr == "" {
		fmt.Println("please specify a valid irc server address")
		os.Exit(1)
	}

	if *username == "" {
		*username = "username"
	}
	if *realname == "" {
		*realname = "realname"
	}
	if *nickname == "" {
		*nickname = "nickname"
	}

	m := manager.New(*log)
	user := fmt.Sprintf("%s 8 * :%s", *username, *realname)
	if err := m.Start(*ssl, *addr, *port, user, *nickname, *password); err != nil {
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := m.Stop(); err != nil {
		os.Exit(1)
	}
}
