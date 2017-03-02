package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/prometheus/common/log"
)

func main() {
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

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("err on starting server-socket on :8080:", err)
	}
	fmt.Println("server running on :8080")

	// CONNECTION TO BROWSER
	client, err := ln.Accept()
	defer client.Close()
	if err != nil {
		log.Fatal("err on client connect", err)
	}
	fmt.Println("client '%v' connected!", client.RemoteAddr())

	fmt.Println("establishing connection to target-server ...")

	// CONNECTION TO NGINX
	target, err := net.Dial("tcp", "127.0.0.1:8181")
	defer target.Close()

	go func() { io.Copy(target, client) }()
	go func() { io.Copy(client, target) }()

	time.Sleep(10 * time.Second)
}
