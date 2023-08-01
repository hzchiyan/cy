package deploy

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "站点部署",
	Long:  "站点部署",
}

func init() {
	Cmd.AddCommand(websiteCmd)
	Cmd.AddCommand(goWebSiteCmd)
}
