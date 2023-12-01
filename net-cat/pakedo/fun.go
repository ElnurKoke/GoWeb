package pakedo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func Message(mess, time string, Client client) message {
	return message{
		name: Client.name,
		text: mess,
		time: time,
	}
}

func HandleClient(conn net.Conn) {
	defer conn.Close()
	art := `
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`
	conn.Write([]byte(art))
	var username string
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	reader := bufio.NewReader(conn)
	for {
		buffer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		buffer = strings.Split(buffer, "\n")[0]
		username = strings.TrimSpace(buffer)
		if len(username) > 0 && len(username) <= 10 && username != "" {
			if _, ok := clients[username]; ok {
				conn.Write([]byte("Already taken, choose another one!\n[ENTER YOUR NAME]:"))
				continue
			} else {
				break
			}
		}
		conn.Write([]byte("Invalid name.[ENTER YOUR NAME]: "))
	}
	// hist, err := os.ReadFile("history.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// if len(string(hist)) > 0 {
	// 	conn.Write([]byte(string(hist) + "\n"))
	// }
	name := username
	Client := client{
		conn: conn,
		// addr: conn.RemoteAddr().String(),
		name: name,
	}

	clientsMux.Lock()
	clients[name] = Client
	joining <- Message("\u001b[38;2;255;160;16m"+"has joined the chat."+"\033[0m", time.Now().Format("2006-01-02 15:04:05"), Client)
	clientsMux.Unlock()
	fmt.Println(fmt.Sprintf("User %s has joined the chat", name))

	// broadcastMessage("\u001b[38;2;255;160;16m"+fmt.Sprintf("User %s has joined the chat", username)+"\033[0m", nil)

	// for {
	// 	// conn.Write([]byte("[" + time.Now().Format("2006-01-02 15:04:05") + "][" + username + "]: "))
	// 	buffer, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error reading:", err)
	// 		break
	// 	}

	// 	message := strings.TrimSpace(buffer)
	// 	if buffer == "" || len(strings.TrimSpace(buffer)) == 0 {
	// 		continue
	// 	}
	// 	now := time.Now().Format("2006-01-02 15:04:05")
	// 	fullMessage := fmt.Sprintf("[%s][%s]: %s", now, username, message)
	// 	fmt.Println(fullMessage)
	// 	AddToHistory(fullMessage)
	// 	broadcastMessage(fullMessage, conn)
	// }
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Text() == "" || len(strings.TrimSpace(input.Text())) == 0 {
			continue
		}
		fullMessage := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), name, input.Text())
		fmt.Println(fullMessage)
		AddToHistory(fullMessage)
		clientsMux.Lock()
		messages <- Message(input.Text(), time.Now().Format("2006-01-02 15:04:05"), Client)
		clientsMux.Unlock()
	}

	clientsMux.Lock()
	delete(clients, name)
	leaving <- Message("\u001b[38;2;255;160;16m"+"has left the chat."+"\033[0m", time.Now().Format("2006-01-02 15:04:05"), Client)
	conn.Close()
	clientsMux.Unlock()

	disconnectMessage := fmt.Sprintf("User %s has left the chat", name)
	fmt.Println(disconnectMessage)
	// broadcastMessage("\u001b[38;2;255;160;16m"+disconnectMessage+"\033[0m", nil)
}

func AddToHistory(message string) error {
	filePath := "history.txt"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	if len(lines) >= 15 {
		lines = lines[1:]
	}

	lines = append(lines, message)

	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ClearFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// func broadcastMessage(message string, sender net.Conn) {
// 	clientsMux.Lock()
// 	defer clientsMux.Unlock()

// 	for client, _ := range clients {
// 		if client != sender {
// 			_, err := client.Write([]byte(message + "\n"))
// 			if err != nil {
// 				fmt.Println("Error writing:", err)
// 			}
// 		} else {
// 			// _, err := client.Write([]byte("\u001b[38;2;0;0;255m" + ("[^You]" + message[:21] + "\n") + "\033[0m"))
// 			// if err != nil {
// 			// 	fmt.Println("Error writing:", err)
// 			// }
// 			continue
// 		}
// 	}
// }
