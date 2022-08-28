package main // import "tls-client"

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var delimiter = '\n'

var id = int(-1)
var fin = make(chan bool)

func sender(con net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("write msg then enter:")

		msgSend, _ := reader.ReadString(byte(delimiter))
		fmt.Fprintf(con, msgSend+string(delimiter))
	}
}

func receiver(con net.Conn) {
	for {
		msgReceive, err := bufio.NewReader(con).ReadString(byte(delimiter))
		if err != nil {
			fmt.Println("connection lost from server" + err.Error())
			fin <- true
		}
		if delimiter == '\n' {
			if strings.Contains(msgReceive, "client id is ") && id == -1 {
				id, err = strconv.Atoi(strings.TrimSpace(strings.Split(msgReceive, " ")[3]))
				if err != nil {
					fmt.Println("error while get id " + err.Error())
					fin <- true
				}
			}

			fmt.Println("from server - " + msgReceive)
		} else {
			fmt.Println("from server no delimeter - " + msgReceive)
		}
	}
}

func tcpClient() {
	port := "8864"

	cert, err := tls.LoadX509KeyPair("key.crt", "key.pem")
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsCfg := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	con, err := tls.Dial("tcp", "127.0.0.1:"+port, tlsCfg)
	if err != nil {
		fmt.Println("error while dialing - " + err.Error())
		return
	}
	defer con.Close()

	fmt.Println("connected. client ready")
	fmt.Println("ctrl+c to finish")

	go sender(con)
	go receiver(con)

	<-fin
}

func main() {
	tcpClient()
}
