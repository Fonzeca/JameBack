/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/Fonzeca/UserHub/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Inicio del server",
	Long:  `Este comando solo inicia la api`,
	Run: func(cmd *cobra.Command, args []string) {
		server.InitServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
