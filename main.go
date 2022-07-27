package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	Conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}

	fmt.Println("WELCOME TO WE-TRANSFER-FILES")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("For registration, enter a username:")
	username, _ := reader.ReadString('\n')

	if len(username) == 0 {
		fmt.Errorf("Enter a valid username")
	}
	macadd, err := getMacAddr()
	if err != nil {
		fmt.Errorf("Something went wrong")
	}

	newClient := Client{
		Con:      Conn,
		username: username,
		macaddr:  macadd,
	}

	fmt.Println("Registration done. \n")
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
