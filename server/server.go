package server

import (
	"fmt"
	"net"
	"os"
        "bufio"
        "strings"
        "time"
)

type MiniserverAisprid struct {
	listener net.Listener
	connection net.Conn
	connected bool
}

func (s MiniserverAisprid) run () {
}

func Listen (port string) (MiniserverAisprid, error) {
	var server = MiniserverAisprid{connected:false}
        listener, err := net.Listen("tcp", port)
        if err != nil {
		fmt.Fprintln(os.Stderr, "Impossible listen on : ",port,".")
		return server, err
        }
        defer listener.Close()
        server.listener = listener

        for {
        	// Waiting connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Impossible accept on port : ",port,".")
			return server, err
		}
		server.connection = conn

        	// Have a connection, read request
                request, err := bufio.NewReader(server.connection).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return server, err
                }
                if strings.TrimSpace(string(request)) == "STOP" {
                        fmt.Println("Exiting TCP server!")
                        return server, nil
                }
                fmt.Print("-> ", string(request))
                
        	// Send response
                t := time.Now()
                myTime := t.Format(time.RFC3339) + "\n"
                server.connection.Write([]byte(myTime))
        }
	return server, nil
}

