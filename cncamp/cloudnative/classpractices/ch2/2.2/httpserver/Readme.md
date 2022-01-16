HTTP 服务器
---

##### Note:

1. 接收 request，并将 request header 信息，写入 response header
2. 读取当前系统环境变量 VERSION 信息，写入 response header
3. Server 侧标准输出访问信息： client IP, HTTP 返回码
4. 当访问 `/healthz` 时，返回 200

#### Dockerfile 

```makefile
BIN_PATH=. REL_VERSION=v1.0 make docker-build

BIN_PATH=. REL_VERSION=v1.0 make docker-push
```