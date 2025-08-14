package main

import (
	"fmt"
	"net"
	"os"
)


// 限制goroutine数量
var limitChan = make(chan bool, 1)

// UDP goroutine 实现并发读取UDP数据
func Process(conn *net.UDPConn)  {
	data := make([]byte, 1024)
	n,remoteAddr,err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Failed To Read UDP Msg, Error: " + err.Error())
	}

	str := string(data[:n])
	fmt.Println("Reveive From Client, Data: " + str)

	server,_ := net.Dial("udp","127.0.0.1:8080")
	server.Write(data[:n])
	go func() {
		defer server.Close()
		for {
			redata := make([]byte, 1024)
			n,_ := server.Read(redata)
			_,err := conn.WriteToUDP(redata[:n],remoteAddr)
			if err != nil{
				fmt.Println(err)
			}
			fmt.Println(server)
			}
	}()
	conn.WriteToUDP([]byte(str), remoteAddr)
	<- limitChan
}

func udpServer(address string)  {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	if err != nil {
		fmt.Println("Read From Connect Failed, Err :" + err.Error())
		os.Exit(1)
	}

	for {
		limitChan <- true
		go Process(conn)
		fmt.Println(limitChan)
	}

}

func main() {
	address := "0.0.0.0:8081"
	udpServer(address)
}