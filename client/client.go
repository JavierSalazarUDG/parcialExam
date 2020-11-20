package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"./message"
)

var client net.Conn
var messages []*message.Message
var nickname string
var show bool = false

func conectClient(nick string) {
	c, err := net.Dial("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Write([]byte(nick))
	err = gob.NewDecoder(c).Decode(&messages)
	if err != nil {
		fmt.Println(err)
	}
	client = c
}
func listenMessages() {
	c := client
	for {
		var message *message.Message
		err := gob.NewDecoder(c).Decode(&message)
		if message.Type > 2 {
			messages = append(messages, message)
			// fmt.Println(messages)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		if show == true {
			if message.Type == 3 {
				fmt.Println(message.User, ":", message.Message)
			} else {
				fmt.Println(message.User, ":", message.FileName)
			}
		}
	}
}
func sendMessage() {
	c := client
	fmt.Println("Escribe un mensaje")
	message := &message.Message{User: nickname, Message: readInpunt(), Type: 3}
	err := gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
	}
}
func readFile() {
	c := client
	var path string
	fmt.Println("Escribe la ruta del archivo")
	fmt.Scanln(&path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	message := &message.Message{User: nickname, Message: "", Type: 4, File: data, FileName: path}
	err = gob.NewEncoder(c).Encode(message)
	if err != nil {
		fmt.Println(err)
	}

}
func showMessages() {
	for _, message := range messages {
		if message.Type == 3 {
			fmt.Println(message.User, ":", message.Message)
		} else {
			fmt.Println(message.User, ":", message.FileName)
		}
	}
}
func readInpunt() string {
	consoleReader := bufio.NewReader(os.Stdin)
	input, err := consoleReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	return input
}
func leave() {
	message := &message.Message{User: nickname, Message: "", Type: 0}
	err := gob.NewEncoder(client).Encode(message)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	loop := true
	//nickname = readInpunt()
	fmt.Scanln(&nickname)
	fmt.Println("hola: ", nickname)
	fmt.Println("conectando al servidor....")
	conectClient(nickname)
	fmt.Println("conectado al servidor")
	go listenMessages()
	defer leave()
	fmt.Println("Escriba la opcion deseada")
	fmt.Println("1.-Enviar mensaje")
	fmt.Println("2.-Enviar un archivo")
	fmt.Println("3.- Mostrar mensajes")
	fmt.Println("4.- Salir")
	for loop == true {
		var opc int
		fmt.Scanln(&opc)
		switch opc {
		case 1:
			sendMessage()
			break
		case 2:
			readFile()
			break
		case 3:
			if show == false {
				showMessages()
			}
			show = !show
			break
		case 4:
			loop = false
			break
		}
	}

}
