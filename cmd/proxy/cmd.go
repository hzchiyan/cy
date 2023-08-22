package proxy

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "proxy",
	Short: "站点部署",
	Long:  "站点部署",
}

func init() {
	Cmd.AddCommand(socket)
}
