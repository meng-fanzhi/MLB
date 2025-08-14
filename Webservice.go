package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"public/loginit"
	"runtime"
)

func Webservice() {
	//第一个参数是接口名，第二个参数 http handle func
	http.HandleFunc("/api/info", Info)
	http.HandleFunc("/api/conn", ConnStat)
	http.HandleFunc("/api/checkheath", CheckHearth)
	//服务器要监听的主机地址和端口号
	http.ListenAndServe("0.0.0.0:8081", nil)
}

type HttpJson struct {
	Errcode int
	Data    interface{}
}

type InfoStruct struct {
	Listen            string
	ListenType        string
	ServerList        string
	ServerType        string
	WhiteList         string
	TimeOut           int
	ClientIdleTimeout int
	ServerIdleTimeout int
	Buffer            int
	Mode              string
	HealthCTime       int
	CurrNumGoroutine  int
}

// http handle func
func ConnStat(rw http.ResponseWriter, req *http.Request) {
	// 返回字符串 "Hello world"
	var data HttpJson
	//dd,_ := json.Marshal()
	data.Data = ServerLists.ServerList
	data.Errcode = 0
	jsondata, _ := json.Marshal(data)
	//json.Marshal()
	_, err := fmt.Fprint(rw, string(jsondata))
	if err != nil {
		loginit.Error.Printf("接口调用失败：%v", err)
	}
}

func Info(rw http.ResponseWriter, req *http.Request) {
	//var data map[string]string
	var data HttpJson
	var info InfoStruct
	info.Listen = LocalIp
	info.ListenType = ListenType
	info.ServerList = ServerAddr
	info.ServerType = ServerType
	info.WhiteList = SourceIps
	info.TimeOut = Timeout
	info.ClientIdleTimeout = ClientIdleTimeout
	info.ServerIdleTimeout = ServerIdleTimeout
	info.Buffer = buffer
	info.Mode = Mode
	HealthCTime := HealthC
	info.CurrNumGoroutine = getNumGroutine()
	if HealthC == 0 || ServerType != "tcp" {
		HealthCTime = 0
	}
	info.HealthCTime = HealthCTime

	data.Data = info
	data.Errcode = 0
	jsondata, _ := json.Marshal(data)
	_, err := fmt.Fprint(rw, string(jsondata))
	if err != nil {
		loginit.Error.Printf("接口调用失败：%v\n", err)
	}
}

func getNumGroutine() int {
	return runtime.NumGoroutine()
}

func CheckHearth(rw http.ResponseWriter, req *http.Request) {
	// 返回字符串 "Hello world"
	//var data HttpJson
	for i := 0; i < len(ServerLists.ServerList); i++ {
		locked.Lock()
		Server.Addr = ServerLists.ServerList[i].Addr
		Server.Count = ServerLists.ServerList[i].Count
		Conn, err := net.Dial(ServerType, Server.Addr)
		if err == nil {
			defer Conn.Close()
			Server.Health = true
			fmt.Printf("# 服务%v(%v)连接正常，健康状态标识为%v,当前连接数为:%v\n", Server.Addr, ServerType, Server.Health, Server.Count)
		} else {
			Server.Health = false
			fmt.Printf("# 服务%v(%v)连接异常，健康状态标识为%v\n", Server.Addr, ServerType, Server.Health)
		}
		ServerLists.ServerList[i] = Server
		locked.Unlock()
	}
	ConnStat(rw, req)
}
