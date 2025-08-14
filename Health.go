package main

import (
	"fmt"
	"net"
	"time"
)

func HealthCheck() {
	for {
		fmt.Println("#########################################################################")
		fmt.Println("# 开启健康状态检查！")
		CheckStat()
		fmt.Println("#########################################################################")
		time.Sleep(time.Duration(HealthC) * time.Second)

	}
}

func CheckStat(){
	//serverlist := ServerLists.ServerList
	for i := 0; i < len(ServerLists.ServerList); i++ {
		go func(i int) {
			locked.Lock()
			Server.Addr = ServerLists.ServerList[i].Addr
			Server.Count = ServerLists.ServerList[i].Count
			Conn, err := net.Dial(ServerType, Server.Addr)
			if err == nil {
				defer Conn.Close()
				Server.Health = true
				fmt.Printf("# 服务%v(%v)连接正常，健康状态标识为%v,当前连接数为:%v\n", Server.Addr,ServerType, Server.Health,Server.Count)
			} else {
				Server.Health = false
				fmt.Printf("# 服务%v(%v)连接异常，健康状态标识为%v\n", Server.Addr,ServerType, Server.Health)
			}
			ServerLists.ServerList[i] = Server
			locked.Unlock()
		}(i)
	}
	//ServerLists.ServerList = serverlist
}
