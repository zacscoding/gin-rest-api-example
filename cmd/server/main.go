package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:  "server",
	Long: "Example of RESTFul Article API Server",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "", "", "config file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("failed to execute command. err: %v", err)
		os.Exit(1)
	}
}
