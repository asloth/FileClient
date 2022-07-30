package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	Con      net.Conn
	username string
}

func (c *Client) Read() error {
	for {
		msg, err := bufio.NewReader(c.Con).ReadBytes('\n')

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
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))

	switch string(cmd) {
	case "RECEIVING":
		c.receiveFile()
	case "SENDING":
		c.sendingFile("notes.txt")
	default:
		fmt.Println(string(cmd))
	}
}

func (c *Client) receiveFile() {
	const BUFFERSIZE = 1024

	connection := c.Con

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")

}

//Function for registering the user in the server
func register(name string, c net.Conn) string {

	command := "REGISTER @" + name + "\n"

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

//Function for sending commands to the server and visualizing the response
func (c *Client) sendCommand(cmd string) {
	_, err := c.Con.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err.Error())
	}
	//Reading the response of the server
	// msg, err := bufio.NewReader(c.Con).ReadBytes('\n')

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(string(msg))
	// return string(msg)
}

//Fuction for listing the channels
func (c *Client) listChannels() {
	command := "LCHANNELS \n"

	c.sendCommand(command)
}

//Function for suscribing to a channel
func (c *Client) suscribing(chann string) {
	command := "SUSCRIBE #" + chann + " \n"
	c.sendCommand(command)

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

	command := "SEND #" + chnn + " \n"

	c.sendCommand(command)

}

//Function for sending only the file data, this is gonna execute at the end of sendFile  method
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
	fmt.Println("File has been sent")
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
