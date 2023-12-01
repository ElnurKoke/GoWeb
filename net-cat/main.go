package main

import (
	"fmt"
	"log"
	"net"
	"net-cat/pakedo"
	"os"
	"strconv"
)

func main() {
	host := "localhost"
	port := 8989
	if len(os.Args) == 2 {
		pickport, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Error port")
			return
		}
		port = pickport
	} else if len(os.Args) > 2 {
		log.Fatal("[USAGE]: ./TCPChat $port")
		return
	}

	pakedo.ClearFile("history.txt")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on %s:%d\n", host, port)
	go pakedo.BroadcastMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
		go pakedo.HandleClient(conn)
	}
}
