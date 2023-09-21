package server

import (
	"fmt"
	"net"
	"os"
        "bufio"
        "strings"
)

type MiniserverAispridClient struct {
	connection net.Conn
	connected bool
}

func Connect (url string) (MiniserverAispridClient, error)  {
	var client = MiniserverAispridClient{connected:false}
        c, err := net.Dial("tcp", url)
        if err != nil {
		fmt.Fprintln(os.Stderr, "Impossible to connect to : ",url,".")
		return client, err
        }
        client.connection = c
        client.connected = true
	fmt.Println("Connected to  : ", url)
	return client, nil
}

func (client MiniserverAispridClient) Get (varName string) (string, error) {
	fmt.Println("Get from server : ", varName)
	return "OK", nil
}

func (client MiniserverAispridClient) Set (varName string, varValue string) (string, error) {
	fmt.Println("Set to server : ", varName, " = ", varValue)
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(client.connection, text+"\n")

        message, _ := bufio.NewReader(client.connection).ReadString('\n')
        fmt.Print("->: " + message)
        if strings.TrimSpace(string(text)) == "STOP" {
                fmt.Println("TCP client exiting...")
		return "OK", nil
        }
	return "OK", nil
}

