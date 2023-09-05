package deploy

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	goWebSiteHost    string
	goWebSiteHostDir string
	goPort           string
)

func init() {
	goWebSiteCmd.Flags().StringVarP(&goWebSiteHost, "host", "H", "", "域名")
	goWebSiteCmd.Flags().StringVarP(&goWebSiteHostDir, "dir", "D", "", "项目部署目录")
	goWebSiteCmd.Flags().StringVarP(&goPort, "port", "P", "8080", "go端口号")
}

var goWebSiteCmd = &cobra.Command{
	Use:   "website-go",
	Short: "部署go-website",
	Long:  "部署go-website",
	Run: func(cmd *cobra.Command, args []string) {
		if goWebSiteHost == "" {
			fmt.Println("--host 参数未传递 域名")
			return
		}
		if goWebSiteHostDir == "" {
			fmt.Println("--dir 参数未传递 域名部署的目录")
			return
		}
		nginxConfig := fmt.Sprintf("server {\n"+
			"listen 80;\n"+
			"root %s;\n"+
			"server_name %s;\n"+
			"location / {\n "+
			"proxy_set_header Host $http_host;\n"+
			"proxy_set_header X-Forwarded-Host $http_host;\n"+
			"proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n"+
			"proxy_set_header X-Real-IP $remote_addr;\n"+
			"proxy_pass http://127.0.0.1:%s;"+
			"\n}\n"+
			"location /.well-known {root %s;}"+
			"\n}\n", goWebSiteHostDir, goWebSiteHost, goPort, goWebSiteHostDir)
		if err := writeFileNginxConfig(goWebSiteHost, nginxConfig); err != nil {
			fmt.Println(err)
			return
		}
		supervisorConfig := fmt.Sprintf("[program:%s]\n"+
			"command=%s/main\n"+
			"directory=%s \n"+
			"autostart=true\n"+
			"autorestart=true\n"+
			"redirect_stderr=true\n"+
			"stdout_logfile=/data/wwwroot/%s.log\n"+
			"stderr_logfile=/data/wwwroot/%s-err.log", goWebSiteHost, goWebSiteHostDir, goWebSiteHostDir, goWebSiteHost, goWebSiteHost)
		if err := writeFilSupervisorConfig(goWebSiteHost, supervisorConfig); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("-------------------------------------")
		fmt.Println("执行一下命令：")
		fmt.Println("nginx -t")
		fmt.Println("nginx -s reload")
		fmt.Println("service supervisor restart")
		fmt.Println("-------------------------------------")
	},
}
