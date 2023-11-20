/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Carmind-Mindia/user-hub/server"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lista de clientes",
	Long: `Muestra una lista de clientes.
	Solo la lista, no se puede mostrar las apikey de cada cliente`,
	Run: func(cmd *cobra.Command, args []string) {
		guard := server.Guard
		list, err := guard.ClientLs()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		for _, v := range list {
			fmt.Println(v)
		}

	},
}

func init() {
	apiKeyCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
