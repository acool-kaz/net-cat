package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	users   = make(map[net.Conn]string)
	msgChan = make(chan Msg, 1)
	allMsg  = []string{}
)

type Msg struct {
	msg    string
	sender string
}

func newMsg(msg, sender string) Msg {
	return Msg{
		msg:    msg,
		sender: sender,
	}
}

func main() {
	port := "8989"
	if len(os.Args[1:]) == 1 {
		port = os.Args[1]
	} else if len(os.Args[1:]) != 0 {
		fmt.Println("[USAGE]: go run server/server.go $port")
		fmt.Println("[USAGE]: go run server/server.go - for starting server on port 8989")
		return
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	if _, err := os.Stat("logs"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	time := time.Now().Format("02 Jan 2006 15:04")
	time = strings.ReplaceAll(time, " ", "_")
	time = strings.ReplaceAll(time, ":", "-")
	file, err := os.Create(fmt.Sprintf("./logs/session-%s.txt", time))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.Printf("Chat started on port :%s\n", port)
	file.WriteString(fmt.Sprintf("Chat started on port :%s\n", port))
	go readMsgFromChan(file)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		printLogo(conn)
		go func() {
			defer conn.Close()
			if err := getName(conn, file); err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			writeMsgToChan(conn)
		}()
	}
}

func printLogo(conn net.Conn) {
	file, err := os.ReadFile("logo.txt")
	if err != nil {
		return
	}
	conn.Write(file)
}

func writeMsgToChan(conn net.Conn) {
	for {
		input := make([]byte, 1024)
		n, err := conn.Read(input)
		if err != nil {
			msgChan <- newMsg(fmt.Sprintf("\n%s has left our chat...\n", users[conn]), users[conn])
			delete(users, conn)
			return
		}
		str := string(input[:n])
		if strings.HasSuffix(str, "\r\n") {
			str = strings.ReplaceAll(str, "\r\n", "\n")
		}
		str = str[:len(str)-1]
		if str == "" {
			conn.Write([]byte(fmt.Sprintf("[%s] [%s]: ", time.Now().Format("2006-01-02 15:04:05"), users[conn])))
			continue
		}
		msgChan <- newMsg(fmt.Sprintf("\n[%s] [%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), users[conn], str), users[conn])
	}
}

func getName(conn net.Conn, file *os.File) error {
	input := make([]byte, 1024)
	name := ""
	for name == "" {
		conn.Write([]byte("[ENTER YOUR NAME]:"))
		n, err := conn.Read(input)
		if err != nil {
			return err
		}
		name = string(input[:n])
		if strings.HasSuffix(name, "\r\n") {
			name = strings.ReplaceAll(name, "\r\n", "\n")
		}
		name = name[:len(name)-1]
	}
	for _, user := range users {
		if user == name {
			return fmt.Errorf("your name is taken, sorry :(")
		}
	}
	if len(users) == 10 {
		return fmt.Errorf("chat limit is 10, sorry :(")
	}
	users[conn] = name
	for _, msg := range allMsg {
		conn.Write([]byte(strings.TrimPrefix(msg, "\n")))
	}
	greetMsg := fmt.Sprintf("\n%s has joined our chat...\n", users[conn])
	file.WriteString(strings.TrimPrefix(greetMsg, "\n"))
	allMsg = append(allMsg, greetMsg)
	for key, user := range users {
		if user == name {
			greetMsg = strings.TrimPrefix(greetMsg, "\n")
		}
		key.Write([]byte(greetMsg))
		key.Write([]byte(fmt.Sprintf("[%s] [%s]: ", time.Now().Format("2006-01-02 15:04:05"), user)))
	}
	return nil
}

func readMsgFromChan(file *os.File) {
	for {
		msg := <-msgChan
		file.WriteString(strings.TrimPrefix(msg.msg, "\n"))
		allMsg = append(allMsg, msg.msg)
		for key, user := range users {
			if user != msg.sender {
				key.Write([]byte(msg.msg))
			}
		}
		for key, user := range users {
			key.Write([]byte(fmt.Sprintf("[%s] [%s]: ", time.Now().Format("2006-01-02 15:04:05"), user)))
		}
	}
}
