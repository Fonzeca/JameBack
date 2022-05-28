/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	guard_userhub "github.com/Fonzeca/UserHub/guard"
	"github.com/Fonzeca/UserHub/server"
	"github.com/Fonzeca/UserHub/server/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "UserHub",
	Short: "Gestor de usuarios",
	Long:  " _ _                _ _       _   \n| | | ___ ___  _ _ | | | _ _ | |_ \n| ' |<_-</ ._>| '_>|   || | || . \\\n`___'/__/\\___.|_|  |_|_|`___||___/\n                                  ",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	utils.InitConfig()

	db, err := server.InitDataBase()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Db = db

	keystore := guard_userhub.NewKeyStore(db) //implements KeyStore interface
	keyGen := guard_userhub.KeyGeneratorUserHub{}

	server.Guard = guard_userhub.NewGuard(&keyGen, keystore)

}
