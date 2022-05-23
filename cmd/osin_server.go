package main

import (
	"avanpost_auth/pkg/osin_server"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "osin_server"}

	rootCmd.AddCommand(osin_server.RootCmdOAuth2Server)
	rootCmd.Execute()
}
