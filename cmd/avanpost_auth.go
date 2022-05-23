package main

import (
	"avanpost_auth/pkg/avanpost_auth"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "avanpost_auth"}
	rootCmd.AddCommand(avanpost_auth.RootCmdAuth)

	err := rootCmd.Execute()
	log.Error().Msg(err.Error())
}
