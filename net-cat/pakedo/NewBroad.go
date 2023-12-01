package pakedo

import (
	"fmt"
	"net"
	"sync"
)

var (
	clientsMux sync.Mutex
	clients    = make(map[string]client) // Mutex для защиты доступа к списку клиентов
	joining    = make(chan message)
	leaving    = make(chan message)
	messages   = make(chan message)
	hst        string
)

type message struct {
	name string
	text string
	time string
}

type client struct {
	conn net.Conn
	name string
}

func BroadcastMessage() {
	for {
		select {
		case msg := <-joining:
			str := "\n" + "\u001b[38;2;255;160;16m" + msg.name + "\033[0m" + " " + msg.text + "\n"
			for _, client := range clients {
				if msg.name != client.name {
					fmt.Fprintf(client.conn, str+"[%s][%s]: ", msg.time, client.name)
				} else {
					fmt.Fprintf(client.conn, hst+"[%s][%s]: ", msg.time, client.name)
				}
			}
		case msg := <-messages:
			str := "\n[" + msg.time + "][" + msg.name + "]: " + msg.text + "\n"
			hst += str[1:]
			for _, client := range clients {
				if msg.name != client.name {
					fmt.Fprintf(client.conn, str)
				}
				fmt.Fprintf(client.conn, "[%s][%s]: ", msg.time, client.name)
			}
		case msg := <-leaving:
			str := "\n" + "\u001b[38;2;255;160;16m" + msg.name + "\033[0m" + " " + msg.text + "\n"
			for _, client := range clients {
				if msg.name != client.name {
					fmt.Fprintf(client.conn, str+"[%s][%s]: ", msg.time, client.name)
				}
			}

		}
	}
}
