> 提供基于4层网络负载均衡转发服务，提供unix_socket,tcp,udp协议的负载均衡服务，并提供基于轮询策略的负载均衡及接入白名单控制等功能；

## 使用说明：
```
Usage of zlb:
-buffer int
-buffer=32 buffer设置,单位为K (default 32)
-cert string
-cert=./cert.pem (default “./certificate/cert.pem”)
-client_idle_timeout int
-client_idle_timeout=0 客户端空闲超时时间,单位为毫秒,0为不限制
-heathcheck int
-heathcheck=300 健康状态检查间隔配置,单位为秒,关闭健康检查配置为0
-idle
-idle 开启空闲超时控制
-key string
-key=./cert.key (default “./certificate/cert.key”)
-limit int
-limit=100000 并发限制,限流 (default 100000)
-listen string
-listen=0.0.0.0:80 指定服务监听的端口 (default “0.0.0.0:8080”)
-listen_type string
-listen_tpye=tcp 监听端协议【tcp,udp】 (default “tcp”)
-mode string
-mode=random 负载均衡模式选择,【polling/random/hash/minc/standby】,polling为轮询算法,random为随机算法,hash算法用于会话保持场景,minc为最小连接数算法,standby为主备切换算法; (default “polling”)
-retry
-retry 开启连接重试
-server string
-server=127.0.0.1:80,多个用’,’隔开，根据策略进行负载均衡 (default “127.0.0.1:80”)
-server_idle_timeout int
-server_idle_timeout=180000 服务端返回空闲超时时间,单位为毫秒,0为不限制 (default 180000)
-server_type string
-server_type=tcp 服务端网络协议【tcp,udp,unix】 (default “tcp”)
-source string
-source=192.168.18.0/24 指定允许访问的源IP网段,多个用’,’隔开 (default “0.0.0.0/0”)
-timeout int
-timeout=30000 客户端空闲超时时间，单位为毫秒 (default 30000)
-tls
-tls 开启TLS传输加密
-udplisten string
-udplisten=0.0.0.0:65534 指定udp转发监听端口，用于接收服务端返回数据 (default “0.0.0.0:65534”)
```

## 举例：
- 将通过转发将docker的unix转换为TCP端口：
```
# 转发docker接口到本地网络2375,限制仅允许x.x.x.x/32访问
mlb -listen=0.0.0.0:2375 -server=//var/run/docker.sock -server_type=unix -src_idle_timeout=0 -dst_idle_timeout=300 -source=x.x.x.x/32
```
```
# 负载均衡后端web服务,并进行链路跟踪
mlb -listen=0.0.0.0:8080 -server=192.168.1.10:8080,192.168.1.11:8080 -server_type=tcp -src_idle_timeout=0 -dst_idle_timeout=300 -mode=hash
```
## QQ讨论群
604869641
