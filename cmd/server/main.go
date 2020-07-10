package main

import (
	"encoding/json"
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
		log.Println(fmt.Sprintf("SecureFShadow 服务器启动成功, 地址: %d, 密码: %d", config.ListenAddr, config.Password))
		log.Println("您的配置文件")
		config_s, _ := json.Marshal(config)
		log.Println(config_s)
	}))
}
