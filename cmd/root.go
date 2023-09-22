
package cmd

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/abriotde/minialertAisprid/server"
)

const (
	EXIT_ARGUMENT_ERROR = 1
)

func runClientCmd(client server.MiniserverAispridClient, args []string) (bool, error) {
	var argsLen = len(args)
	if argsLen<1 {
		fmt.Fprintln(os.Stderr, "You must give comands to client.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	} else if argsLen<2 {
		fmt.Fprintln(os.Stderr, "Missing arguments.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
	var clientCmd = args[0]
	if clientCmd=="send" {
		if argsLen<3 {
			fmt.Fprintln(os.Stderr, "Missing arguments.")
			os.Exit(EXIT_ARGUMENT_ERROR)
		}
		var varName = args[1]
		varValue,err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Argument 2 must be an integer for the value : ",args[2],".")
			return false, err
		}
		// TODO : check varname/varvalue match possible value (No injection)
        	fmt.Println("Send to server : ", varName, " = ", varValue)
        	_,err = client.Set(varName, int32(varValue))
        	if err!=nil {
        		return false, err
        	}
	} else if clientCmd=="get" {
		var varName = args[1]
		// TODO : check varname match possible value (No injection)
        	fmt.Println("Get from server : ", varName)
        	_,err := client.Get(varName);
        	if err!=nil {
        		return false, err
        	}
	} else {
		fmt.Fprintln(os.Stderr, "Unknown client command : ", clientCmd, " possibilities are send|get.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
        return true, nil
}

var (
	// Used for flags.
	interactive     bool
	serverURL		string
	port		string
	rootCmd = &cobra.Command{
		Use:   "minialertAisprid",
		Short: "MinialertAisprid is a minimalistic chalenge to send messages and receive alerts.",
		Long: `Minialert is a minimalistic chalenge consisting in a client/server which send messages and receive alerts lists
			This is licenced under GPL V3`,
		Run: func(cmd *cobra.Command, args []string) {
			if port!="" {
				// var server MiniserverAisprid
				fmt.Println("Run as server mode on port ", port, ".")
				_, err := server.Listen("localhost:"+port)
				if err != nil {
					fmt.Fprintln(os.Stderr, " : ",err,".")
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
			} else {
				var client, err = server.Connect(serverURL)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Fail to connect to server : ",serverURL,".")
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
				if interactive {
			    		fmt.Println("Interactive mode is enable.")
			    	} else {
			    		runClientCmd(client, args)
					fmt.Println("Client: " + strings.Join(args, " "))
			   	}
			   	client.Close()
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
}

// --interactive -i
// --server -s --client -c
// send battery|cpu XX
// get alerts
func init() {
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "For client, it run on interactive mode.")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "localhost:8080", "The server to connect when in client mode (default). If no port specified, it connect on 8080 port." )
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "The port to connect so run it as server.")
}

