package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//Connecting to the server
	Conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	//Welcoming the user
	fmt.Println("WELCOME TO WE-TRANSFER-FILES")

	//Setting the reader
	reader := bufio.NewReader(os.Stdin)

	//Asking for a username
	fmt.Println("For registration, enter a username:")

	//Creating the client
	var newClient Client

	//Reading the username until it gets the OK message in the register method
	for {
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		//Validating that the username is not empty
		if len(username) == 0 {
			fmt.Println("Enter a valid username")
			continue
		}
		if len(username) > 10 {
			fmt.Println("A username can not be longer than 10 digits")
			continue
		}
		//Registering the user
		msg := register(username, Conn)
		msg = strings.TrimSpace(msg)
		//If msg is OK then break the loop and show the menu
		if msg == "OK" {
			newClient = Client{
				Con:      Conn,
				username: username,
			}

			break
		}

		fmt.Println(msg)
	}

	fmt.Println("Registration done.")
	fmt.Println("Menu: \n 1. Listar canales \n 2. Subscribise a un canal \n 3. Enviar archivo a un canal \n 4. Abandonar un canal \n 5. Salir")
	go newClient.Read()

menu:
	for {
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			err := newClient.listChannels()
			if err != nil {
				fmt.Println(err.Error(), " Returning to menu.")
				continue
			}

		case "2":
			fmt.Println("Enter the name of the channel you want to join:")
			channelName, _ := reader.ReadString('\n')
			channelName = strings.TrimSpace(channelName)

			err := newClient.suscribing(channelName)
			if err != nil {
				fmt.Println(err.Error(), " Returning to menu.")
				continue
			}
		case "3":
			fmt.Println("Enter the name of the channel you want sent the file to:")
			channelName, _ := reader.ReadString('\n')
			channelName = strings.TrimSpace(channelName)
			fmt.Println("Enter the path to the file you want to send:")
			filePath, _ := reader.ReadString('\n')
			filePath = strings.TrimSpace(filePath)

			newClient.sendFile(channelName, filePath)
		case "4":
			fmt.Println("Enter the name of the channel you want to left:")
			channelName, _ := reader.ReadString('\n')
			channelName = strings.TrimSpace(channelName)
			err := newClient.unsuscribing(channelName)
			if err != nil {
				fmt.Println(err.Error(), " Returning to menu.")
				continue
			}

		case "5":
			fmt.Println("Good bye")
			newClient.Con.Close()
			break menu
		default:
			fmt.Println("Enter a number from 1 to 4 please")
		}
	}

}
