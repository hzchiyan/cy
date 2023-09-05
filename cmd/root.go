package cmd

import (
	"github.com/hzchiyan/cy/cmd/deploy"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cy",
	Short: "cy",
	Long:  `cy`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deploy.Cmd)
}
