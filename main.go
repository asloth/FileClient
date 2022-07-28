package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	Conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}

	fmt.Println("WELCOME TO WE-TRANSFER-FILES")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("For registration, enter a username:")
	var newClient Client

	for {
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if len(username) == 0 {
			fmt.Println("Enter a valid username")
			continue
		}

		msg := register(username, Conn)
		msg = strings.TrimSpace(msg)

		if msg == "OK" {
			newClient = Client{
				Con:      Conn,
				username: username,
			}

			break
		}
	}

	fmt.Println("Registration done.")
	fmt.Println("Menu: \n 1. Listar canales \n 2. Subscribise a un canal \n 3. Enviar archivo a un canal")
	option, _ := reader.ReadString('\n')

	switch option {
	case "1":
		fmt.Println("im 1")
	case "2":
		newClient.receiveFile()
	case "3":
		fmt.Println("im 3")
	}

}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
