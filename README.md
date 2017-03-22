# websocket_tester
stress test of websocket 

# System Setting

for support more connections, need to change system parameters

* sysctl -w fs.file-max=10485760 #系统允许的文件描述符数量10m
* sysctl -w net.ipv4.tcp_rmem=1024 #每个tcp连接的读取缓冲区1k，一个连接1k，10m只需要10G
* sysctl -w net.ipv4.tcp_wmem=1024 #每个tcp连接的写入缓冲区1k
* echo '* soft nofile 1048576' >> /etc/security/limits.conf #用户单进程的最大文件数，用户登录时生效
* echo '* hard nofile 1048576' >> /etc/security/limits.conf #用户单进程的最大文件数，用户登录时生效
* ulimit -n 1048576 #用户单进程的最大文件数 当前会话生效

# Usage 

```shell
main -cfg config.json
```

# config.json

JSON format

```json
{
  "ws_scheme":"ws",
  "ws_ip":"192.168.33.11",
  "ws_port": "7272",
  "ws_path": "/",
  "str_login":"{\"type\":\"login\",\"user_id\":\"%d\",\"event_id\":\"1\",\"group_id\":\"1\",\"client_name\":\"1\",\"room_id\":\"1\"}",
  "str_say":"{\"type\":\"say\",\"parameters\":{\"body\":\"維持 session ${USER_ID}\",\"title_id\":\"10\",\"user_id\":\"${USER_ID}\",\"nick_name\":\"u${USER_ID}\",\"event_id\":\"1\",\"group_id\":\"1\",\"type\":\"1\",\"icon\":\"icon.png\"},\"to_client_id\":\"1\",\"content\":\"維持 session  ${USER_ID}\"}",
  "str_ping":"{\"type\":\"ping\"}",
  "str_pong":"{\"type\":\"pong\"}",
  "simulator_count":2,
  "simulator_start_in":5,
  "exec_second": 12
}
```

**ws_scheme** websocket protocal, **ws** or **wss**
**ws_ip** websocket server's IP/domain name
**ws_port** the public port
**ws_path** websocket path
**str_login** the login type message
**str_say** the say type message
**str_ping** server's ping message
**str_pong** response for server's ping message
**simulator_count** start how many client to connect to server
**simulator_start_in** start all client in how many time(second)
**exec_second** how long to run


