/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Carmind-Mindia/user-hub/server"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Crea un cliente y muestra la apikey",
	Long: `Crea un cliente y muestra la apikey
La apikey solo se muestra una vez.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		guard := server.Guard
		key, err := guard.GenerateAndSaveApiKey(args[0])
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fmt.Println("ApiKey:")
		fmt.Println("-----------------------------------------")
		fmt.Println(key.ReadableValue())
		fmt.Println("-----------------------------------------")
		fmt.Println("Solo se muestra una vez.")
	},
}

func init() {
	apiKeyCmd.AddCommand(newCmd)
}
