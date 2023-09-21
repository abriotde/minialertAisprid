package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type MiniserverAisprid struct {
    serverName string
    serverPort  int
}

func (s MiniserverAisprid) run () {
}

type MiniserverAispridClient struct {
    serverHostname string
    serverPort  int
    connected bool
}

func Connect (url string) (MiniserverAispridClient, error)  {
	host, port, err := net.SplitHostPort(url)
	// urlinfos, err := url.Parse(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fail to parse url to server : ",url,".")
		return MiniserverAispridClient{}, err
	}
	portValue, err := strconv.Atoi(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "port is not a valid integer : ",port,".")
		return MiniserverAispridClient{}, err
	}
	var client = MiniserverAispridClient{serverHostname:host, serverPort:portValue, connected:true}
	fmt.Println("Connected to  : ", host," on port ", portValue)
	return client, nil;
}

func (s MiniserverAispridClient) Get (varName string) {
	fmt.Println("Get from server : ", varName)
}

func (s MiniserverAispridClient) Set (varName string, varValue string) {
	fmt.Println("Set to server : ", varName, " = ", varValue)
}

