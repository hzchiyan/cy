package proxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
)

var (
	socketAddr     string
	socketUsername string
	socketPassword string
)

func init() {
	socket.Flags().StringVarP(&socketAddr, "addr", "A", "0.0.0.0:7890", "地址")
	socket.Flags().StringVarP(&socketUsername, "username", "U", "", "用户名")
	socket.Flags().StringVarP(&socketPassword, "password", "P", "8080", "密码")
}

var socket = &cobra.Command{
	Use:   "socket",
	Short: "启动socket服务",
	Long:  "启动socket服务",
	Run: func(cmd *cobra.Command, args []string) {
		listener, err := net.Listen("tcp", socketAddr)
		if err != nil {
			log.Fatalln("服务启动失败:", err)
		}
		log.Println("服务启动成功：", socketAddr)
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("接受失败:", err)
				continue
			}
			go handleSocketConnection(conn, socketUsername, socketPassword)
		}
	},
}

func handleSocketConnection(conn net.Conn, username string, password string) {
	// 读取客户端发送的版本信息
	buffer := make([]byte, 256)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println("读取版本错误:", err)
		return
	}
	// 确认版本信息为Socket5
	if buffer[0] != 0x05 {
		log.Println("不支持的socket5版本")
		return
	}
	if username == "" || password == "" {
		// 发送支持的认证方法
		conn.Write([]byte{0x05, 0x00})
	} else {
		conn.Write([]byte{0x05, 0x02})
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
		}
		// 读取账号密码
		ulen := int(buffer[1])
		u := string(buffer[2 : 2+ulen])
		plen := int(buffer[2+ulen])
		p := string(buffer[3+ulen : 3+ulen+plen])
		if username == u && password == p {
			conn.Write([]byte{0x05, 0x00})
		} else {
			conn.Write([]byte{0x05, 0x01})
			return
		}
	}
	// 读取客户端发送的请求信息
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("读取请求失败:", err)
		return
	}
	// 确认请求为连接命令
	if buffer[1] != 0x01 {
		log.Println("不支持的连接类型")
		return
	}
	// 解析请求地址
	var address string
	switch buffer[3] {
	case 0x01: // IPv4地址
		address = fmt.Sprintf("%d.%d.%d.%d:%d", buffer[4], buffer[5], buffer[6], buffer[7], uint16(buffer[n-2])<<8|uint16(buffer[n-1]))
	case 0x03: // 域名
		address = string(buffer[5:n-2]) + ":" + fmt.Sprintf("%d", uint16(buffer[n-2])<<8|uint16(buffer[n-1]))
	default:
		log.Println("错误的地址")
		return
	}
	// 连接目标服务器
	targetConn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("远程服务器连接失败:", err)
		return
	}
	log.Println("连接成功：", conn.RemoteAddr().String(), "-->", targetConn.RemoteAddr().String())
	// defer targetConn.Close()
	// 发送响应给客户端，表示连接已建立
	response := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	conn.Write(response)
	// 进行转发数据
	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)
}
