/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manejo de apikeys",
	Long: `Apikeys que maneja userhub
Los clientes utilizan las apikeys para consumir el servicio de userhub.`,
}

func init() {
	rootCmd.AddCommand(apiKeyCmd)
	apiKeyCmd.SetUsageTemplate(
		`Usage:
  UserHub apikey new [client-name]`)
}
