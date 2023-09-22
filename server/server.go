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
        "google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc"
	"github.com/abriotde/minialertAisprid/messages"
	"github.com/abriotde/minialertAisprid/monitorer"
)

type MiniserverAisprid struct {
	listener net.Listener
	connection net.Conn
	connected bool
}

var monitoring monitorer.Monitorer

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

// 
func (s *server_t) SendDataMetric(ctx context.Context, in *messages.SendDataMetricRequest) (*messages.SendDataMetricReply, error) {
	sValue := strconv.Itoa(int(in.GetValue()))
	fmt.Println("Received: ", in.GetName(), " = ", sValue)
	monitoring.Log(in.GetName(), in.GetValue())
	return &messages.SendDataMetricReply{Message: "Set " + in.GetName() + " = "+sValue, Ok:true}, nil
}
// 
func (s *server_t) GetAlertHistory(ctx context.Context, in *messages.GetAlertHistoryRequest) (*messages.GetAlertHistoryReply, error) {
	fmt.Println("Ask for alerts.")
	var alerts []*messages.GetAlertHistoryReply_Alert
	var nbAlerts = 0
	for _, alert := range monitoring.GetAlertHistory() {
		a := messages.GetAlertHistoryReply_Alert{Timestamp:timestamppb.New(alert.Timestamp), Name:alert.Name, Value:alert.Value}
		alerts = append(alerts, &a)
		nbAlerts += 1
	}
	fmt.Println("Have ",strconv.Itoa(nbAlerts)," alerts.")
	return &messages.GetAlertHistoryReply{Alerts:alerts, Ok:true}, nil
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

