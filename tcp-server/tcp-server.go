package main // import "tcp-server"

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	// only needed below for sample processing
)

var delimiter = '\n'
var cons = make(map[string]net.Conn)

func sendMessage(con net.Conn, msg string) {
	fmt.Fprintf(con, msg+string(delimiter))
}

func sender() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("write msg then enter: \n")
		msgToSend, _ := reader.ReadString(byte(delimiter))

		if len(cons) == 0 {
			fmt.Println("no clients connected")
			continue
		}

		msgs := strings.Split(msgToSend, " ")
		if len(msgs) > 1 {
			targetID := msgs[0]
			message := msgs[1]

			if cons[targetID] != nil {
				sendMessage(cons[targetID], message)
				continue
			}

			_, err := strconv.ParseInt(targetID, 10, 64)
			if err == nil {
				fmt.Println("no connection with id " + targetID)
				continue
			}
		}

		// Broadcast
		for _, con := range cons {
			sendMessage(con, msgToSend)
		}
	}
}

func receiver(con net.Conn, id string) {
	fmt.Println("client id is " + id)
	_, err := con.Write([]byte("client id is " + id + "\n"))
	if err != nil {
		fmt.Println("error while writing " + err.Error())
	}

	for {
		msg, _ := bufio.NewReader(con).ReadString(byte(delimiter))

		if len(msg) > 0 {
			if delimiter == '\n' {
				fmt.Print("from client:", string(msg))
			} else {
				fmt.Println("from client no delimeter:", string(msg))
			}

			newMsg := "echo~ " + strings.ToUpper(strings.TrimSuffix(msg, string(delimiter)))
			con.Write([]byte(newMsg + string(delimiter)))
		} else {
			fmt.Println("closed from client " + id)

			delete(cons, id)
			con.Close()
			break
		}
	}
}

func tcpServer() {
	fmt.Println("Server Listening TCP localhost:7749")
	fmt.Println("Ctrl+C to finish.")

	listener, _ := net.Listen("tcp", ":7749")
	defer listener.Close()

	go sender()

	for {
		con, err := listener.Accept()
		if err != nil {
			panic("error while accepting " + err.Error())
		}

		id := fmt.Sprint(rand.Intn(100))
		fmt.Println("connection accepted")

		cons[id] = con
		go receiver(con, id)
	}
}

func main() {
	tcpServer()
}
