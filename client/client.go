package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ip := "localhost"
	port := "8989"
	if len(os.Args[1:]) == 2 {
		ip = os.Args[1]
		port = os.Args[2]
	} else if len(os.Args[1:]) != 0 {
		fmt.Println("[USAGE]: ./TCPChat $ip $port")
		fmt.Println("[USAGE]: ./TCPChat - for connection to localhost:8989")
		return
	}
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		for {
			input := make([]byte, 1024)
			n, err := conn.Read(input)
			if err != nil {
				os.Exit(0)
			}
			fmt.Print(string(input[:n]))
		}
	}()
	for {
		str, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(str))
	}
}
