package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/xiangrui2019/securefshadow/config"
	"github.com/xiangrui2019/securefshadow/local"
)

func main() {
	config := &config.Config{
		ListenAddr: ":1800",
	}

	config.ReadConfig()
	config.SaveConfig()

	local, err := local.NewLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatalln(local.Listen(func(listenAddr *net.TCPAddr) {
		log.Println("SecureFShadow 本地Socks5端启动成功.")
		log.Println("您的配置文件: ")
		config_s, _ := json.Marshal(config)
		log.Println(string(config_s))
	}))
}
