package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"./clientList"
	"./message"
)

var clients []clientList.Client
var messages []*message.Message
var flag = true

func server() {
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := server.Accept()
		b := make([]byte, 100)
		bs, err := c.Read(b)
		err = gob.NewEncoder(c).Encode(messages)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go listenMessages(c)
		clients = append(clients, clientList.Client{Client: c, Nickname: string(b[:bs])})
		fmt.Println("New client conected: ", string(b[:bs]))
	}
}
func sendMessages(message *message.Message) {
	if message.Type < 3 {
		return
	}
	for _, client := range clients {
		err := gob.NewEncoder(client.Client).Encode(&message)
		if err != nil {
			deleteClient(client.Nickname)
			continue
		}
	}
}
func listenMessages(client net.Conn) {
	for {
		var message *message.Message
		err := gob.NewDecoder(client).Decode(&message)
		if err != nil {
			break
		}
		if message.Type == 0 {
			deleteClient(message.User)
		}
		if message.Type > 2 {
			messages = append(messages, message)
			sendMessages(message)
		}

	}
}
func deleteClient(nickname string) {
	for index, c := range clients {
		if nickname == c.Nickname {
			fmt.Println("Removing client: ", nickname)
			clients = RemoveIndex(clients, index)
		}
	}
}
func RemoveIndex(s []clientList.Client, index int) []clientList.Client {
	return append(s[:index], s[index+1:]...)
}
func showMessages() {
	if messages == nil {
		return
	}
	for _, message := range messages {
		if message.Type == 3 {
			fmt.Println(message.User, ":", message.Message)
		} else {
			fmt.Println(message.User, ":", message.FileName)
		}
	}
}

func downloadMessages() {
	f, err := os.Create("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, message := range messages {
		var line string
		if message.Type == 3 {
			line = message.User + ":" + message.Message
			fmt.Println(message.User, ":", message.Message)
		} else {
			line = message.User + ":" + message.FileName
			fmt.Println(message.User, ":", message.FileName)
		}
		_, err := f.WriteString(line)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func serverMenu() {

	for flag == true {
		var opc int
		fmt.Println("Escriba la opcion deseada")
		fmt.Println("1.-Mostrar mensajes")
		fmt.Println("2.-Respaldar mensajes")
		fmt.Println("3.- Terminar servidor")
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			showMessages()
			break
		case 2:
			downloadMessages()
			break
		case 3:
			flag = false
			break
		default:
			//loop = false
		}
	}
}
func main() {
	go server()
	fmt.Println("Server listen on port 3000")
	go serverMenu()
	for flag == true {

	}

}
