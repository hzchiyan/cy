package deploy

import (
	"github.com/hzchiyan/cy/internal/fs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "站点部署",
	Long:  "站点部署",
}

func init() {
	Cmd.AddCommand(websiteCmd)
	Cmd.AddCommand(goWebSiteCmd)
	Cmd.AddCommand(htmlWebSiteCmd)
}

func writeFileNginxConfig(host, nginxConfig string) error {
	nginxfile := "/etc/nginx/conf.d/" + host + ".conf"
	if _, err := fs.WriteFile(nginxfile, nginxConfig); err != nil {
		return err
	}
	return nil
}

func writeFilSupervisorConfig(host, nginxConfig string) error {
	nginxfile := "/etc/supervisor/conf.d/" + host + ".conf"
	if _, err := fs.WriteFile(nginxfile, nginxConfig); err != nil {
		return err
	}
	return nil
}
