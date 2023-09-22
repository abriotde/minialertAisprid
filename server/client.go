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

// It is client for connecting to "MiniserverAisprid" server.
type MiniserverAispridClient struct {
	connection *grpc.ClientConn 
	grpcConnection messages.GreeterClient
	connected bool
}

// Fonction to connect to the server : prerequisites for all communications : call to GetAlertHistory and Set
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
// Fonction to close connection to the server
func (client MiniserverAispridClient) Close ()  {
	client.connection.Close()
}

// Return alert history based on previous Set() call where name/value exceed threashold defined on server side.
func (client MiniserverAispridClient) GetAlertHistory () ([]*messages.GetAlertHistoryReply_Alert, error) {
	fmt.Println("GetAlertHistory from server : ")

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.GetAlertHistory(ctx, &messages.GetAlertHistoryRequest{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not GetAlertHistory", err)
		return make([]*messages.GetAlertHistoryReply_Alert, 0), err
	}
	if r.GetOk() != true {
		fmt.Fprintln(os.Stderr, "could not GetAlertHistory : Server refuse.")
		return make([]*messages.GetAlertHistoryReply_Alert, 0), errors.New("could not GetAlertHistory : Server refuse.")
	}
	return r.GetAlertHistory(), nil
}
// Set the current value of the variable name registred on the server.
func (client MiniserverAispridClient) Set (varName string, varValue int32) (string, error) {
	fmt.Println("Set to server : ", varName, " = ", varValue)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.SendDataMetric(ctx, &messages.SendDataMetricRequest{Name:varName, Value:varValue})
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not SendDataMetric : ",varName," = ",varValue,": ", err)
	}
	if r.GetOk() != true {
		fmt.Fprintln(os.Stderr, "could not SendDataMetric : ",varName," = ",varValue,": Server refuse.")
	}
	fmt.Println("Greeting: ", r.GetMessage())
	return "OK", nil
}
// Function to test client/server communication with simple ASCII (useless now?)
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

