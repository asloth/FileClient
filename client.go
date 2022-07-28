package main

import (
	"bufio"
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

func (c *Client) err(e error) {
	c.Con.Write([]byte("ERR " + e.Error() + "\n"))
}

func (c *Client) read() error {
	for {
		msg, err := bufio.NewReader(c.Con).ReadBytes('\n')
		if err == io.EOF {
			// Connection closed, deregister client
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println(string(msg))
	}
}
