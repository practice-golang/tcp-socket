package main // import "tcp-server"

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	// only needed below for sample processing
)

var delimiter = '\n'
var Msg chan string = make(chan string, 10)
var Health chan string = make(chan string)

func pcListener(fin chan int) {
	fmt.Println("Server Listening TCP this IP:7749")
	fmt.Println("Ctrl+C to finish.")

	ln, _ := net.Listen("tcp", ":7749")
	defer ln.Close()

	con, _ := ln.Accept()
	defer con.Close()

	for {
		select {
		case request := <-Msg:
			con.Write([]byte(request + string(delimiter)))
		case health := <-Health:
			healthMsg := strings.ToUpper(strings.TrimSuffix(health, string(delimiter)))
			con.Write([]byte(healthMsg + string(delimiter)))
			msg, _ := bufio.NewReader(con).ReadString(byte(delimiter))
			if len(msg) > 0 {
				if delimiter == '\n' {
					fmt.Print("Msg Received:", string(msg))
				} else {
					fmt.Println("Msg Received:", string(msg))
				}

				// newMsg := strings.ToUpper(strings.TrimSuffix(msg, string(delimiter)))
			} else {
				fmt.Println("Closed from client.")
				con.Close()
				con, _ = ln.Accept()
			}
		}
	}
	// fin <- 1
}

func main() {
	fin := make(chan int)
	go pcListener(fin)
	// go func() {
	i := 0
	for {
		log.Println("hahaha" + fmt.Sprint(i))
		Msg <- "hahaha" + fmt.Sprint(i)
		time.Sleep(time.Second * 1)
		if i%2 == 0 {
			Health <- "health"
		}
		i++
	}
	// }()
	// <-fin
}
