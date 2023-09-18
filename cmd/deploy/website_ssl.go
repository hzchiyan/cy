package deploy

import (
	"fmt"
	"github.com/hzchiyan/cy/internal/fs"
	"github.com/hzchiyan/cy/internal/ids"
	"github.com/hzchiyan/cy/internal/model"
	"github.com/hzchiyan/cy/internal/ssl"
	"github.com/spf13/cobra"
)

var (
	host string
	dir  string
)

func init() {
	websiteCmd.Flags().StringVarP(&host, "host", "H", "", "域名")
	websiteCmd.Flags().StringVarP(&dir, "dir", "D", "", "项目部署目录")
}

var websiteCmd = &cobra.Command{
	Use:   "website-ssl",
	Short: "部署站点ssl",
	Long:  "部署站点ssl",
	Run: func(cmd *cobra.Command, args []string) {
		email := ids.MacID() + "@hzchiyangithub.com"
		if host == "" {
			fmt.Println("--host 参数未传递 域名")
			return
		}
		if dir == "" {
			fmt.Println("--dir 参数未传递 域名部署的目录")
			return
		}
		websiteSSL, err := ssl.WebsiteSSL(model.DB, email, host, dir)
		if err != nil {
			fmt.Println(fmt.Sprintf("client.ObtainSSL err=%v", err))

			fmt.Println("---------------nginx config append-------------------")
			fmt.Println("可能需要修改nginx配置,追加下面这行指令")
			fmt.Println(fmt.Sprintf("location /.well-known {root %s;}", dir))
			fmt.Println("---------------nginx config-------------------")

			return
		}
		keyFile := "/etc/nginx/" + host + ".key"
		pemFile := "/etc/nginx/" + host + ".pem"
		_, err = fs.WriteFile(keyFile, websiteSSL.PrivateKey)
		if err != nil {
			fmt.Println(keyFile + "写入失败")
			return
		}
		_, err = fs.WriteFile(pemFile, websiteSSL.Pem)
		if err != nil {
			fmt.Println(pemFile + "写入失败")
			return
		}
		fmt.Println("---------------nginx config append-------------------")
		fmt.Println(fmt.Sprintf(""+
			"listen 443 http2; \n"+
			"ssl_certificate  %s;\n  "+
			"ssl_certificate_key  %s;\n  "+
			"ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;\n    "+
			"ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;\n   "+
			"ssl_prefer_server_ciphers on;\n    "+
			"ssl_session_cache shared:SSL:10m;\n   "+
			"ssl_session_timeout 10m;", pemFile, keyFile))
		fmt.Println("---------------nginx config-------------------")
	},
}
