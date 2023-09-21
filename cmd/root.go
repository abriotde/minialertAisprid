
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
  Use:   "minialert",
  Short: "Minialert is a minimalistic chalenge to send messages and receive alerts.",
  Long: `Minialert is a minimalistic chalenge consisting in a client/server which send messages and receive alerts lists
          This is licenced under GPL V3`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

