package main

import (
	"hash/crc32"
	"math/rand"
	"os"
	"public/loginit"
	"strings"
)

type ServerStruct struct {
	Addr   string
	Health bool
	Count  int
}

type ServerListStruct struct {
	ServerList []ServerStruct
}

var Server ServerStruct
var ServerLists ServerListStruct

func InitServer() {
	SourceIp = strings.Split(SourceIps, ",")
	ServerAddrs := strings.Split(ServerAddr, ",")

	for i := 0; i < len(ServerAddrs); i++ {
		Server.Addr = ServerAddrs[i]
		//dconn, err := net.Dial(dtype, ServerAddrs[i])
		Server.Health = true
		Server.Count = 0
		ServerLists.ServerList = append(ServerLists.ServerList, Server)
	}

	//启动健康状态检查服务
	if HealthC != 0 && ServerType == "tcp" {
		go HealthCheck()
	}

	if len(ServerAddrs) <= 0 {
		loginit.Error.Println("服务IP和端口不能空,或者无效")
		os.Exit(1)
	}

	switch Mode {
	case "polling":
		//var LbModeInter Balancer.Balancer
		LbInter = Polling{}
	case "random":
		LbInter = Random{}
	case "hash":
		LbInter = Hash{}
	case "minc":
		LbInter = MinConn{}
	case "standby":
		LbInter = Standby{}
	default:
		loginit.Error.Println("未识别负载均衡模式，请确认后重新配置！")
		os.Exit(1)
	}
}

// 负载均衡接口
type Balancer interface {
	StartLb(string) (string, bool)
}

// 轮询方法
type Polling struct {
}

func (Polling) StartLb(string) (string, bool) {
	locked.Lock()
	defer locked.Unlock()
	//fmt.Printf("aaa:%v\n",ServerLists.ServerList)
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Health {
			ip := ServerLists.ServerList[i].Addr
			ServerLists.ServerList = append(ServerLists.ServerList[i+1:], ServerLists.ServerList[:i+1]...)
			return ip, true
		}
	}
	return "", false
}

// 随机方法
type Random struct {
}

func (Random) StartLb(string) (string, bool) {
	locked.Lock()
	defer locked.Unlock()
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Health {
			lens := len(ServerLists.ServerList)
			index := rand.Intn(lens)
			ip := ServerLists.ServerList[index].Addr
			return ip, true
		}
	}
	return "", false
}

// 一致性哈希算法
type Hash struct {
}

func (Hash) StartLb(remote string) (string, bool) {
	locked.Lock()
	defer locked.Unlock()

	defKey := remote
	crcTable := crc32.MakeTable(crc32.IEEE)
	hashVal := crc32.Checksum([]byte(defKey), crcTable)

	lens := len(ServerLists.ServerList)
	index := int(hashVal) % lens
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[(index+i)%lens].Health {
			ip := ServerLists.ServerList[(index+i)%lens].Addr
			return ip, true
		}
	}
	return "", false
}

// 最小连接数法
type MinConn struct {
}

func (MinConn) StartLb(string) (string, bool) {
	locked.Lock()
	defer locked.Unlock()
	MinCount := 65535
	ip := ""
	//fmt.Printf("%v",ServerLists.ServerList)
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Health && ServerLists.ServerList[i].Count < MinCount {
			MinCount = ServerLists.ServerList[i].Count
			ip = ServerLists.ServerList[i].Addr
		}
	}
	if MinCount != 65535 && ip != "" {
		return ip, true
	}
	return "", false
}

// 主备切换
type Standby struct {
}

func (Standby) StartLb(string) (string, bool) {
	locked.Lock()
	defer locked.Unlock()
	//fmt.Printf("%v",ServerLists.ServerList)
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Health {
			ip := ServerLists.ServerList[i].Addr
			return ip, true
		}
	}
	return "", false
}
