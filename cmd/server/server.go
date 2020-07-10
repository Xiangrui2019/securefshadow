package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/xiangrui2019/securefshadow/config"
	"github.com/xiangrui2019/securefshadow/lib/password"
	"github.com/xiangrui2019/securefshadow/server"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("SERVER_LISTEN_PORT"))

	if err != nil {
		port = 7748
	}

	config := &config.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		Password:   password.RandPassword(),
	}

	config.ReadConfig()
	config.SaveConfig()

	server, err := server.NewServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
SecureFShadow服务器 启动成功, 配置文件如下:
服务监听地址:
%s
密码
%s`, listenAddr, config.Password))
	}))
}
