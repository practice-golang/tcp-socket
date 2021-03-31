package main // import "tcp-client"

import (
	"bufio"
	"fmt"
	"net"
)

var delimiter = '\n'

func main() {
	con, _ := net.Dial("tcp", "127.0.0.1:7749")
	fmt.Println("Client ready.")
	fmt.Println("Ctrl+C to finish.")
	for {
		// reader := bufio.NewReader(os.Stdin)

		// fmt.Print("Msg for sending then press Enter: ")

		// msgSend, _ := reader.ReadString(byte(delimiter))
		// fmt.Fprintf(con, msgSend+string(delimiter))

		msgReceive, err := bufio.NewReader(con).ReadString(byte(delimiter))
		if err != nil {
			fmt.Println("Connection lost.")
			panic(err)
		}
		if delimiter == '\n' {
			fmt.Fprintf(con, msgReceive+string(delimiter)) // 차임벨클라 헬스 응답
			fmt.Print("Response from server: " + msgReceive)
		} else {
			fmt.Println("Response from server: " + msgReceive)
		}
	}
}
