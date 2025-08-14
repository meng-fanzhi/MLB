package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// 限制goroutine数量
var limitChan = make(chan bool, 10)

// UDP goroutine 实现并发读取UDP数据
func Process(conn *net.UDPConn,i int)  {
	fmt.Println(i)
	data := make([]byte, 1024)
	n,remoteAddr,err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Failed To Read UDP Msg, Error: " + err.Error())
	}

	str := string(data[:n])
	fmt.Println("Reveive From Client, Data: " + str)
	ii := 0
	for{
		time.Sleep(1000)
		conn.WriteToUDP([]byte("Reveive From Client, Data: " + str + "\n"), remoteAddr)
		ii++
		if ii == 100 {
			break
		}
	}



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
	i := 0
	for {
		limitChan <- true
		i++
		go Process(conn,i)
		fmt.Println(limitChan)
	}

}

func main() {
	address := "0.0.0.0:8080"
	udpServer(address)
}