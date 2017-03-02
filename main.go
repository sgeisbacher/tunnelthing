package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	target string
	client string
	port   int
)

func init() {
	flag.StringVar(&target, "target", "", "the target (<host>:<port>)")
	flag.StringVar(&target, "t", "", "the target (<host>:<port>) (shorthand)")
	flag.IntVar(&port, "port", 7757, "the tunnelthing port")
	flag.IntVar(&port, "p", 7757, "the tunnelthing port (shorthand)")
}

func main() {
	flag.Parse()
	if target == "" {
		log.Fatal("no target specified")
	}

	signals := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for _ = range signals {
			fmt.Println("\nReceived an interrupt, stopping services...")
			stop <- true
		}
	}()

	printWelcome()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("err on starting server-socket on %d: %v", port, err)
	}
	fmt.Printf("server running on %d\n", port)

	// CONNECTION TO BROWSER
	client, err := ln.Accept()
	defer client.Close()
	if err != nil {
		log.Fatal("err on client connect", err)
	}
	fmt.Printf("client '%v' connected!\n", client.RemoteAddr())

	fmt.Println("establishing connection to target-server ...")
	// CONNECTION TO NGINX
	target, err := net.Dial("tcp", target)
	fmt.Printf("connection to server %v established!\n", target.RemoteAddr())
	defer target.Close()

	go func() { io.Copy(target, client) }()
	go func() { io.Copy(client, target) }()

	<-stop
}

func printWelcome() {
	fmt.Println(`
 __          __  _                            _          _                          _ _   _     _             _ 
 \ \        / / | |                          | |        | |                        | | | | |   (_)           | |
  \ \  /\  / /__| | ___ ___  _ __ ___   ___  | |_ ___   | |_ _   _ _ __  _ __   ___| | |_| |__  _ _ __   __ _| |
   \ \/  \/ / _ \ |/ __/ _ \| '_   _ \ / _ \ | __/ _ \  | __| | | | '_ \| '_ \ / _ \ | __| '_ \| | '_ \ / _  | |
    \  /\  /  __/ | (_| (_) | | | | | |  __/ | || (_) | | |_| |_| | | | | | | |  __/ | |_| | | | | | | | (_| |_|
     \/  \/ \___|_|\___\___/|_| |_| |_|\___|  \__\___/   \__|\__,_|_| |_|_| |_|\___|_|\__|_| |_|_|_| |_|\__, (_)
                                                                                                         __/ |  
                                                                                                        |___/   	
	`)
}
