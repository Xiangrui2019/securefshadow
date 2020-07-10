package local

import (
	"log"
	"net"

	"github.com/xiangrui2019/securefshadow/lib/encryption"
	"github.com/xiangrui2019/securefshadow/lib/password"
	"github.com/xiangrui2019/securefshadow/lib/transport"
)

type Local struct {
	Cipher     *encryption.Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func NewLocal(password_i string, listenAddr, remoteAddr string) (*Local, error) {
	bsPassword, err := password.ParsePassword(password_i)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	structRemoteAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &Local{
		Cipher:     encryption.NewCipher(bsPassword),
		ListenAddr: structListenAddr,
		RemoteAddr: structRemoteAddr,
	}, nil
}

// 本地端启动监听，接收来自本机浏览器的连接
func (local *Local) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return transport.ListenSecureTCP(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}

func (local *Local) handleConn(userConn *transport.SecureTCPConnection) {
	defer userConn.Close()

	proxyServer, err := transport.DialSecureTCP(local.RemoteAddr, local.Cipher)
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()

	// 进行转发
	// 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := proxyServer.DecodeCopy(userConn)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	userConn.EncodeCopy(proxyServer)
}
