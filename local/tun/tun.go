package tun

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/xiangrui2019/securefshadow/local"
)

const listenAddr = "127.0.0.1:9400"

func StartTunServer(password string, remoteAddr string, fd int) (err error) {
	lsLocal, err := local.NewLocal(password, listenAddr, remoteAddr)
	if err != nil {
		return
	}
	lwipStack := core.NewLWIPStack()
	f := os.NewFile(uintptr(fd), "tun")
	if f == nil {
		return errors.New("无法打开VPN虚拟网卡")
	}
	defer f.Close()
	lsLocal.Listen(func(listenAddr *net.TCPAddr) {
		log.Println("StartTunServer", listenAddr)
		// 成功启动LS服务器，开启转发
		core.RegisterTCPConnHandler(socks.NewTCPHandler(listenAddr.IP.String(), uint16(listenAddr.Port)))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(listenAddr.IP.String(), uint16(listenAddr.Port), 2*time.Minute))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			return f.Write(data) // 写入tun
		})
		io.Copy(lwipStack, f) // 从tun中读
	})
	return
}
