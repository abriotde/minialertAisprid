package server

import (
	"fmt"
	"net"
	"os"
        "time"
        "context"
        "strconv"
        "bufio"
        "strings"
	"google.golang.org/grpc"
	"github.com/abriotde/minialertAisprid/messages"
)

type MiniserverAisprid struct {
	listener net.Listener
	connection net.Conn
	connected bool
}

// server is used to implement helloworld.GreeterServer.
type server_t struct {
	messages.UnimplementedGreeterServer
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
        server.Run()

	return server, nil
}

// SayHello implements helloworld.GreeterServer
func (s *server_t) SetIntVar(ctx context.Context, in *messages.SetIntVarRequest) (*messages.SetIntVarReply, error) {
	sValue := strconv.Itoa(int(in.GetValue()))
	fmt.Println("Received: ", in.GetName(), " = ", sValue)
	return &messages.SetIntVarReply{Message: "Set " + in.GetName() + " = "+sValue, Ok:true}, nil
}

func (server MiniserverAisprid) Run () (MiniserverAisprid, error) {
	grpcServer := grpc.NewServer()
	messages.RegisterGreeterServer(grpcServer, &server_t{})
	fmt.Println("server listening at ", server.listener.Addr())
	if err := grpcServer.Serve(server.listener); err != nil {
		fmt.Fprintln(os.Stderr, "failed to serve: %v", err)
		return server, err
	}
	return server, nil
}

func (server MiniserverAisprid) Test () (MiniserverAisprid, error) {
        for {
        	// Waiting connection
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Impossible accept.")
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

