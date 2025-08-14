package main

import (
	"fmt"
	"strconv"
	"strings"
)

// 子网计算
func subnet(source_ip string, subnet_net string) bool {
	var maskLen int
	var addrInt2 int64
	addrInt := ipAddrToInt(source_ip)

	addrInt2arry := strings.Split(subnet_net, "/")

	if len(addrInt2arry) != 2 {
		maskLen = 32
		addrInt2 = ipAddrToInt(subnet_net)
	} else {
		maskLen, _ = strconv.Atoi(addrInt2arry[1])
		addrInt2 = ipAddrToInt(addrInt2arry[0])
	}

	fmt.Printf("%v", maskLen)
	/* 十进制转化为二进制 */
	//c := strconv.FormatInt(addrInt, 2)
	//fmt.Println("c:", c)
	addrInt_mask := strconv.FormatInt(addrInt>>uint(32-maskLen), 2)
	addrInt2_mask := strconv.FormatInt(addrInt2>>uint(32-maskLen), 2)
	return addrInt_mask == addrInt2_mask
}

// 网络地址转化
func ipAddrToInt(ipAddr string) int64 {
	bits := strings.Split(ipAddr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}
