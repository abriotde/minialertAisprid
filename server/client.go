package server

import (
	"fmt"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/abriotde/minialertAisprid/messages"
	"context"
	"time"
	"errors"
)

type MiniserverAispridClient struct {
	connection *grpc.ClientConn
	grpcConnection messages.GreeterClient
	connected bool
}

func Connect (url string) (MiniserverAispridClient, error)  {
	var client = MiniserverAispridClient{connected:false}
        // conn, err := net.Dial("tcp", url)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
		fmt.Fprintln(os.Stderr, "Impossible to connect to : ",url,".")
		return client, err
        }
        client.connection = conn
        client.connected = true
	client.grpcConnection = messages.NewGreeterClient(client.connection)
	fmt.Println("Connected to  : ", url)
	return client, nil
}
func (client MiniserverAispridClient) Close ()  {
	client.connection.Close()
}

func (client MiniserverAispridClient) GetAlerts () ([]*messages.GetAlertsReply_Alert, error) {
	fmt.Println("GetAlerts from server : ")

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.GetAlerts(ctx, &messages.GetAlertsRequest{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not GetAlerts", err)
		return make([]*messages.GetAlertsReply_Alert, 0), err
	}
	if r.GetOk() != true {
		fmt.Fprintln(os.Stderr, "could not GetAlerts : Server refuse.")
		return make([]*messages.GetAlertsReply_Alert, 0), errors.New("could not GetAlerts : Server refuse.")
	}
	return r.GetAlerts(), nil
}

func (client MiniserverAispridClient) Set (varName string, varValue int32) (string, error) {
	fmt.Println("Set to server : ", varName, " = ", varValue)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.SetIntVar(ctx, &messages.SetIntVarRequest{Name:varName, Value:varValue})
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not SetIntVar : ",varName," = ",varValue,": ", err)
	}
	if r.GetOk() != true {
		fmt.Fprintln(os.Stderr, "could not SetIntVar : ",varName," = ",varValue,": Server refuse.")
	}
	fmt.Println("Greeting: ", r.GetMessage())
	return "OK", nil
}

func (client MiniserverAispridClient) Test () (string, error) {
	/* fmt.Println("Test mode : client")
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(client.connection, text+"\n")

        message, _ := bufio.NewReader(client.connection).ReadString('\n')
        fmt.Print("->: " + message)
        if strings.TrimSpace(string(text)) == "STOP" {
                fmt.Println("TCP client exiting...")
		return "OK", nil
        } */
	return "OK", nil
}

