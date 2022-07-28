package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

//Function for sending commands to the server and visualizing the response
func (c *Client) sendCommand(cmd string) {
	_, err := c.Con.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err.Error())
	}
	//Reading the response of the server
	msg, err := bufio.NewReader(c.Con).ReadBytes('\n')

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(msg))
}

//Fuction for listing the channels
func (c *Client) listChannels() {
	command := "LCHANNELS \n"

	c.sendCommand(command)
}

func (c *Client) suscribing(chann string) {
	command := "SUSCRIBE #" + chann + " \n"
	c.sendCommand(command)

	go c.receiveFile()
}

func (c *Client) sendFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	pr, pw := io.Pipe()
	w, err := gzip.NewWriterLevel(pw, 7)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		n, err := io.Copy(w, file)
		if err != nil {
			log.Fatal(err)
		}
		w.Close()
		pw.Close()
		log.Printf("copied to piped writer via the compressed writer: %d", n)
	}()

	n, err := io.Copy(c.Con, pr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("copied to connection: %d", n)
}
