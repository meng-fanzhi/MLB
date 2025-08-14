package main

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"
)

var locked sync.Mutex
var SourceIps string
var SourceIp []string
var RemoteIp []string
var LocalIp string
var UdpListener string
var TLS bool
var TlsCert string
var TlsKey string
var Limit int
var ListenType string
var ServerType string
var ServerAddr string
var ServerIdleTimeout int
var ClientIdleTimeout int
var Timeout int
var buffer int
var Mode string
var HealthC int
var ReTry bool
var Idle bool
var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")
var LbInter Balancer
var Counting Conut
var HttpApi bool

func main() {
	flag.StringVar(&LocalIp, "listen", "0.0.0.0:8080", "-listen=0.0.0.0:80 指定服务监听的端口")
	flag.StringVar(&UdpListener, "udplisten", "0.0.0.0:65534", "-udplisten=0.0.0.0:65534 指定udp转发监听端口，用于接收服务端返回数据")
	flag.BoolVar(&TLS, "tls", false, "-tls 开启TLS传输加密")
	flag.StringVar(&TlsCert, "cert", "./certificate/cert.pem", "-cert=./cert.pem")
	flag.StringVar(&TlsKey, "key", "./certificate/cert.key", "-key=./cert.key")
	flag.StringVar(&ServerAddr, "server", "127.0.0.1:80", "-server=127.0.0.1:80,多个用','隔开，根据策略进行负载均衡")
	flag.StringVar(&SourceIps, "source", "0.0.0.0/0", "-source=192.168.18.0/24 指定允许访问的源IP网段,多个用','隔开")
	flag.StringVar(&ListenType, "listen_type", "tcp", "-listen_tpye=tcp 监听端协议【tcp,udp】")
	flag.StringVar(&ServerType, "server_type", "tcp", "-server_type=tcp 服务端网络协议【tcp,udp,unix】")
	flag.IntVar(&ClientIdleTimeout, "client_idle_timeout", 0, "-client_idle_timeout=0 客户端空闲超时时间,单位为毫秒,0为不限制")
	flag.IntVar(&ServerIdleTimeout, "server_idle_timeout", 180000, "-server_idle_timeout=180000 服务端返回空闲超时时间,单位为毫秒,0为不限制")
	flag.IntVar(&Timeout, "timeout", 30000, "-timeout=30000 客户端空闲超时时间，单位为毫秒")
	flag.BoolVar(&ReTry, "retry", false, "-retry 开启连接重试")
	flag.BoolVar(&Idle, "idle", false, "-idle 开启空闲超时控制")
	flag.IntVar(&buffer, "buffer", 32, "-buffer=32 buffer设置,单位为K")
	flag.IntVar(&Limit, "limit", 100000, "-limit=100000 并发限制,限流")
	flag.StringVar(&Mode, "mode", "polling", "-mode=random 负载均衡模式选择,【polling/random/hash/minc/standby】,polling为轮询算法,random为随机算法,hash算法用于会话保持场景,minc为最小连接数算法,standby为主备切换算法;")
	flag.IntVar(&HealthC, "heathcheck", 0, "-heathcheck=300 健康状态检查间隔配置,单位为秒,关闭健康检查配置为0")
	flag.BoolVar(&HttpApi, "http", false, "-http 开启http查询接口")
	flag.Parse()

	//初始化服务列表
	InitServer()
	fmt.Println("#########################################################################")
	fmt.Println("# 初始化完成！")
	fmt.Println("#########################################################################")
	time.Sleep(time.Duration(1) * time.Second)
	if HttpApi {
		go Webservice()
	}

	//启动服务
	server()
}
