package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"public/loginit"
	"strconv"
	"strings"
	"time"
)


// 监听服务
func server() {
	fmt.Println("#########################################################################")
	fmt.Println("#  欢迎使用智云转发平台-Deadline V4.2.1     ")
	fmt.Println("#  Dev BY LaoMeng    ")
	fmt.Println("#########################################################################")
	fmt.Printf("#  监听：%v(%v)\n", LocalIp, ListenType)
	fmt.Printf("#  服务地址：%v(%v)\n", ServerAddr, ServerType)
	fmt.Printf("#  白名单：%v\n", SourceIps)
	fmt.Printf("#  并发限制为：%v \n", Limit)
	fmt.Printf("#  连接超时：%v 毫秒\n", Timeout)

	if Idle {
		fmt.Printf("#  服务端空闲超时：%v 毫秒\n", ServerIdleTimeout)
		fmt.Printf("#  客户端空闲超时：%v 毫秒\n", ClientIdleTimeout)
		fmt.Printf("#  buffer:%v K\n", buffer)
	}

	fmt.Printf("#  TLS加密传输:%v\n", TLS)

	if TLS {
		fmt.Printf("#  证书文件：%v\n", TlsCert)
		fmt.Printf("#  私钥文件：%v\n", TlsKey)
	}

	fmt.Printf("#  连接重试：%v\n", ReTry)
	fmt.Printf("#  空闲超时断开：%v\n", Idle)
	fmt.Printf("#  负载均衡模式：%v\n", Mode)
	fmt.Printf("#  开启HTTP API：%v\n", HttpApi)
	HealthCTime := HealthC
	if HealthC == 0 || ServerType != "tcp" {
		HealthCTime = 0
	}
	fmt.Printf("#  健康状态监测间隔：%v 秒\n", HealthCTime)
	fmt.Printf("#  %v\n", time.Now())
	fmt.Println("#########################################################################")

	// 限制goroutine数量

	if ListenType == "udp" {
		handleUDP(LocalIp)

	} else {
		if TLS {
			cer, err := tls.LoadX509KeyPair(TlsCert, TlsKey)
			if err != nil {
				loginit.Error.Println(err)
				return
			}
			tlscon := &tls.Config{Certificates: []tls.Certificate{cer}}
			listls, err := tls.Listen(ListenType, LocalIp, tlscon)
			if err != nil {
				loginit.Error.Println(err)
				return
			}
			defer listls.Close()
			AcceptListen(listls)

		} else {
			lis, err := net.Listen(ListenType, LocalIp)
			if err != nil {
				loginit.Error.Println(err)
				return
			}
			defer lis.Close()
			AcceptListen(lis)
		}
	}
}

func AcceptListen(lis net.Listener) {
	var limitChan = make(chan bool, Limit)
	for {
		conn, err := lis.Accept()
		if err != nil {
			loginit.Error.Printf("建立连接错误:%v\n", err)
			continue
		}
		RemoteIp = strings.Split(conn.RemoteAddr().String(), ":")
		stat := false
		Counting = Conut{}

		for i := 0; i < len(SourceIp); i++ {
			if subnet(RemoteIp[0], SourceIp[i]) {
				stat = true
				break
			}
		}
		if stat {
			//fmt.Println("源地址:%v\n，目的地址源地址:%v", conn.RemoteAddr().String(), conn.LocalAddr().String())
			ip, ok := LbInter.StartLb(RemoteIp[0])
			if !ok {
				loginit.Error.Printf("获取后端服务地址异常，获取地址为:%v\n", ip)
				conn.Close()
				continue
			} else {
				loginit.Info.Printf("源地址:%v,目的地址:%v\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
				go handleTCP(conn, ip, ReTry,&limitChan)
			}
		} else {
			loginit.Warning.Printf("连接被拒绝:%v\n", conn.RemoteAddr().String())
			conn.Close()
			continue
		}
	}
}

// 建立TCP链接请求
func handleTCP(sconn net.Conn, ip string, Retry bool,limitChan *chan bool) {
	*limitChan <- true
	defer sconn.Close()
	//remote_ip = strings.Split(sconn.RemoteAddr().String(), ":")
	dconn, err := net.DialTimeout(ServerType, ip, time.Duration(Timeout)*time.Millisecond)
	if err != nil {
		loginit.Error.Printf("连接服务%v失败,失败原因:%v\n", ip, err)
		if Retry {
			loginit.Warning.Printf("连接重试,连接服务%v\n", ip)
			handleTCP(sconn, ip, false,limitChan)
		}
		if HealthC != 0 {
			CheckStat()
		}
		return
	}
	Counting.Add(ip)
	defer Counting.Del(ip)
	defer dconn.Close()

	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		var dwritten int64
		if Idle {
			loginit.Info.Printf("开启客户端空闲超时：%v", ClientIdleTimeout)
			dwritten, err = copyBuffer(dconn, sconn, ClientIdleTimeout, nil)
		} else {
			dwritten, err = io.CopyBuffer(dconn, sconn, nil)
		}
		Filesize := formatFileSize(dwritten)
		loginit.Info.Printf("服务%v断开连接,断开原因:%v,接收数据量%v\n", ip, err, Filesize)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		var swritten int64
		if Idle {
			loginit.Info.Printf("开启服务端空闲超时：%v", ServerIdleTimeout)
			swritten, err = copyBuffer(sconn, dconn, ServerIdleTimeout, nil)
		} else {
			swritten, err = io.CopyBuffer(sconn, dconn, nil)
		}
		Filesize := formatFileSize(swritten)
		loginit.Info.Printf("客户端%v断开连接,断开原因:%v,接收数据量%v\n", sconn.RemoteAddr(), err, Filesize)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	<-*limitChan
	dconn.Close()
}

func handleUDP(address string){
	var limitChan = make(chan bool, Limit)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	listener, err := net.ListenUDP("udp", udpAddr)
	udpAddr1, err := net.ResolveUDPAddr("udp", UdpListener)
	listener1, err := net.ListenUDP("udp", udpAddr1)

	defer listener.Close()
	if err != nil {
		fmt.Println("Read From Connect Failed, Err :" + err.Error())
		os.Exit(1)
	}
	//i := 0
	for {
		//i++
		//print(i)
		limitChan <- true
		go udpServer(listener,listener1,&limitChan)
	}
}

func udpServer(listener *net.UDPConn,listener1 *net.UDPConn,limitChan *chan bool) {
	data := make([]byte, buffer * 1024)
	n, remoteAddr, err := listener.ReadFromUDP(data)
	//loginit.Info.Printf(listener.LocalAddr().String())
	if err != nil {
		loginit.Error.Printf("Failed To Read UDP Msg, Error: " + err.Error())
	}

	remote_ip := strings.Split(remoteAddr.String(), ":")

	// 连接服务端
	ip, ok := LbInter.StartLb(remote_ip[0])
	if !ok {
		return
	}
	rAddr, _ := net.ResolveUDPAddr("udp", ip)
	//写入数据
	_, conn_write_err := listener1.WriteToUDP(data[:n],rAddr)
	if conn_write_err != nil {
		loginit.Warning.Printf("Write To UDP Server Error: " + conn_write_err.Error())
	}

	for {
		ideaTime, _ := time.ParseDuration(strconv.Itoa(ServerIdleTimeout) + "ms")
		err = listener1.SetReadDeadline(time.Now().Add(ideaTime))
		if err != nil{
			loginit.Info.Printf(err.Error())
		}

		redata := make([]byte, buffer * 1024)
		n,_, Read_err := listener1.ReadFromUDP(redata)
		if Read_err != nil {
			loginit.Warning.Printf("Read From UDP Server Error: " + Read_err.Error())
			break
		}
		loginit.Info.Printf("写入%v",redata[:n] )
		_, err := listener.WriteToUDP(redata[:n],remoteAddr)
		if err != nil {
			loginit.Info.Printf("返回数据失败%v",err )
			break
		}
	}
	<- *limitChan
}
