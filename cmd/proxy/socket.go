package proxy

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
)

var (
	socketAddr string
)

func init() {
	socket.Flags().StringVarP(&socketAddr, "addr", "A", "0.0.0.0:1080", "地址")
}

var socket = &cobra.Command{
	Use:   "socket",
	Short: "启动socket服务",
	Long:  "启动socket服务",
	Run: func(cmd *cobra.Command, args []string) {
		s := Socket{Addr: socketAddr}
		s.Start()
	},
}

type Socket struct {
	Addr, User, Passwd string
}

func (s *Socket) Start() {
	server, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return
	}
	log.Println("服务启动成功：", server.Addr())
	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v", err)
			continue
		}
		go s.callback(client)
	}
}
func (s *Socket) callback(conn net.Conn) {
	if err := s.noAuth(conn); err != nil {
		_ = conn.Close()
		log.Printf("socket auth: err =%v", err)
		return
	}
	target, err := s.connect(conn)
	if err != nil {
		target.Close()
		log.Printf("socket connect: err =%v", err)
		return
	}
	log.Println("conn：", conn.RemoteAddr().String(), "-->", target.RemoteAddr().String())
	s.forward(conn, target)
}
func (s *Socket) noAuth(client net.Conn) (err error) {
	buf := make([]byte, 256)
	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}
	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}
	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}
	// 通知客户端无需认证
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp err: " + err.Error())
	}
	return nil
}
func (s *Socket) auth(client net.Conn) (err error) {
	buf := make([]byte, 256)
	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}
	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}
	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}
	//     无认证的socks5
	if buf[0] == 0x00 {
		n, err = client.Write([]byte{0x05, 0x00})
		if n != 2 || err != nil {
			return errors.New("write rsp err: " + err.Error())
		}
		// 允许无认证
		//return nil
		return errors.New("no auth")
	}
	// 带认证的socks5
	if buf[0] == 0x02 {
		n, err = client.Write([]byte{0x05, 0x02})
		if n != 2 {
			return errors.New("reading methods: " + err.Error())
		}
		// 检查请求头版本
		n, err = io.ReadFull(client, buf[:1])
		//认证子协商版本（与 SOCKS 协议版本的0x05无关系）
		if buf[0] != 0x01 {
			return errors.New("reading methods: " + err.Error())
		}
		//从请求中获取用户名
		n, err = io.ReadFull(client, buf[:1])
		userLen := buf[0]
		n, err = io.ReadFull(client, buf[:userLen])
		bufUsername := string(buf[:userLen])
		n, err = io.ReadFull(client, buf[:1])
		passwdLen := buf[0]
		n, err = io.ReadFull(client, buf[:passwdLen])
		bufPassword := string(buf[:passwdLen])
		if bufUsername == s.User && bufPassword == s.Passwd {
			//STATUS：认证结果（0x00 认证成功 / 大于0x00 认证失败）
			n, err = client.Write([]byte{0x05, 0x00})
			return nil
		} else {
			n, err = client.Write([]byte{0x05, 0x01})
			return errors.New("auth error")
		}
	}
	return nil
}
func (s *Socket) connect(conn net.Conn) (net.Conn, error) {
	// 定义一个数组，用来xxx
	buf := make([]byte, 256)
	n, err := io.ReadFull(conn, buf[:4])
	// 检查是否为ipv4协议
	if n != 4 {
		return nil, errors.New("read header: " + err.Error())
	}
	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}
	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(conn, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])
		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	case 4:
		return nil, errors.New("IPv6: no supported yet")
	default:
		return nil, errors.New("invalid atyp")
	}
	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])
	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, errors.New("dial dst: " + err.Error())
	}
	n, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}
	return dest, nil
}

func (s *Socket) forward(conn, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(conn, target)
	go forward(target, conn)
}
