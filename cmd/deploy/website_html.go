package deploy

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	htmlWebSiteHost    string
	htmlWebSiteHostDir string
)

func init() {
	htmlWebSiteCmd.Flags().StringVarP(&htmlWebSiteHost, "host", "H", "", "域名")
	htmlWebSiteCmd.Flags().StringVarP(&htmlWebSiteHostDir, "dir", "D", "", "项目部署目录")
}

var htmlWebSiteCmd = &cobra.Command{
	Use:   "website-html",
	Short: "部署html-website",
	Long:  "部署html-website",
	Run: func(cmd *cobra.Command, args []string) {
		if htmlWebSiteHost == "" {
			fmt.Println("--host 参数未传递 域名")
			return
		}
		if htmlWebSiteHostDir == "" {
			fmt.Println("--dir 参数未传递 域名部署的目录")
			return
		}
		nginxConfig := fmt.Sprintf("server {\n"+
			"listen 80;\n"+
			"root %s;\n"+
			"server_name %s;\n"+
			"location /.well-known {root %s;}"+
			"\n}\n", htmlWebSiteHostDir, htmlWebSiteHost, htmlWebSiteHostDir)
		if err := writeFileNginxConfig(htmlWebSiteHost, nginxConfig); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("-------------------------------------")
		fmt.Println("执行nginx -t")
		fmt.Println("执行nginx -s reload")
		fmt.Println("-------------------------------------")
	},
}
