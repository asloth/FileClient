package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Client struct {
	Con      net.Conn
	username string
}

func (c *Client) Read() error {
	for {
		msg := make([]byte, 2)
		_, err := c.Con.Read(msg)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		c.Handle(msg)
	}
}

func (c *Client) Handle(message []byte) {
	//Taking the command from the received message
	cmd := bytes.ToUpper(bytes.TrimSpace(message))

	switch string(cmd) {
	case "RC":
		c.receiveFile()
	case "OK":
		fmt.Println("OK")
	default:
		msg := make([]byte, 20)
		c.Con.Read(msg)
		fmt.Println(string(cmd) + string(msg))
	}
}

// Function for registering the user in the server
func register(name string, c net.Conn) string {
	//Completando el nombre hasta los 10 bytes requeridos
	fullName := fillString(name, 10)

	//Defining the command that is gonna be send to the server
	command := "REG@" + fullName

	//Writing the command in the connection
	_, err := c.Write([]byte(command))

	if err != nil {
		return err.Error()
	}
	//Reading the response of the server
	msg, err := bufio.NewReader(c).ReadBytes('\n')

	if err != nil {
		return err.Error()
	}

	return string(msg)
}

// Function for sending commands to the server and visualizing the response
func (c *Client) sendCommand(cmd string) {
	_, err := c.Con.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Fuction for listing the channels
func (c *Client) listChannels() error {
	command := "LCH"

	c.sendCommand(command)
	return nil
}

// Function for suscribing to a channel
func (c *Client) suscribing(chann string) error {

	//Validating that the username is not empty
	if len(chann) == 0 {
		return fmt.Errorf("enter a valid channel name")
	}
	if len(chann) > 10 {
		return fmt.Errorf("a channel name can not be longer than 10 digits")
	}

	//Completando el nombre hasta los 10 bytes requeridos
	channelName := fillString(chann, 10)

	command := "SUS#" + channelName
	c.sendCommand(command)
	return nil

}

// Function for unsuscribing to a channel
func (c *Client) unsuscribing(chann string) error {

	//Validating that the username is not empty
	if len(chann) == 0 {
		return fmt.Errorf("enter a valid channel name")
	}
	if len(chann) > 10 {
		return fmt.Errorf("a channel name can not be longer than 10 digits")
	}

	//Completando el nombre hasta los 10 bytes requeridos
	channelName := fillString(chann, 10)

	command := "UNS#" + channelName

	c.sendCommand(command)
	return nil

}

// function for sending a file to a channel
func (c *Client) sendFile(chnn, path string) {
	//checking if the file exists
	_, error := os.Stat(path)

	// check if error is "file not exists"
	if os.IsNotExist(error) {
		fmt.Printf("%v file does not exist. Returning to the menu\n", path)
		return
	}

	command := "SND#" + fillString(chnn, 10)

	c.sendCommand(command)

	c.sendingFile(path)
}

// Function for sending only the file data, this is gonna execute at the end of sendFile  method
func (c *Client) sendingFile(path string) {
	const BUFFERSIZE = 1024

	connection := c.Con
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}

}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func (c *Client) receiveFile() {
	const BUFFERSIZE = 1024

	connection := c.Con

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	_, err := connection.Read(bufferFileSize)
	if err != nil {
		panic(err)
	}

	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	_, err = connection.Read(bufferFileName)

	if err != nil {
		panic(err)
	}

	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := create("received/" + fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	fmt.Println("Start receiving the file")
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")

}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
