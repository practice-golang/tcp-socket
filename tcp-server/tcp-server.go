package main // import "tcp-server"

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	// only needed below for sample processing
)

var delimiter = '\n'

func pcListener(fin chan int) {
	fmt.Println("Server Listening TCP localhost:8081")
	fmt.Println("Ctrl+C to finish.")

	ln, _ := net.Listen("tcp", ":8081")
	defer ln.Close()

	con, _ := ln.Accept()
	defer con.Close()

	for {
		msg, _ := bufio.NewReader(con).ReadString(byte(delimiter))

		if len(msg) > 0 {
			if delimiter == '\n' {
				fmt.Print("Msg Received:", string(msg))
			} else {
				fmt.Println("Msg Received:", string(msg))
			}

			newMsg := strings.ToUpper(strings.TrimSuffix(msg, string(delimiter)))
			con.Write([]byte(newMsg + string(delimiter)))
		} else {
			fmt.Println("Closed from client.")
			con.Close()
			con, _ = ln.Accept()
		}
	}
	fin <- 1
}

func main() {
	fin := make(chan int)
	go pcListener(fin)
	<-fin
}
