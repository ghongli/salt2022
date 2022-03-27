Istio Ingress Gateway 的形式发布 httpserver
---

### 准备

#### httpserver 服务镜像

`service0` 接收请求后，向 `service1` 发送请求，同时收集请求头 `request header` 如下的 key，并传播到其他请求：

```wiki
x-request-id
x-b3-traceid
x-b3-spanid
x-b3-parentspanid
x-b3-sampled
x-b3-flags
x-ot-span-context
```

##### 构建 `service0` docker 镜像

```shell
cd service0
BIN_PATH=. REL_VERSION=v1.0 make docker-build

BIN_PATH=. REL_VERSION=v1.0 make docker-push
```

##### 构建 `service1` docker 镜像

```shell
cd service1
BIN_PATH=. REL_VERSION=v1.0 make docker-build

BIN_PATH=. REL_VERSION=v1.0 make docker-push
```

### k8s 中部署 `service0`、`service1`

#### 创建 `tracing` 命令空间，开启 istio 自动注入功能

```shell
kubectl create ns tracing
kubectl label ns tracing istio-injection=enabled
```

#### 部署 `service0`、`service1`

```shell
kubectl -n tracing apply -f service0/k8s/deployment.yaml
kubectl -n tracing apply -f service1/k8s/deployment.yaml
```

#### 查看部署 pod

```shell
kubectl -n tracing get po
```

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0omx3njzaj20iy02tgm3.jpg" alt="tracing-po" style="zoom:80%;" />

### Istio 暴露 service

#### `VirtualService` 中配置路由：`/http/service0` -> `/service0`

```yaml
http:
  - match:
      - uri:
          exact: '/http/service0'
    rewrite:
      uri: '/service0'
    route:
      - destination:
          host: service0
            port:
              number: 80
          weight: 100
    timeout: 15s
    retries:
      attempts: 3
      perTryTimeout: 5s
```

#### 创建 TLS 证书和私钥

```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 -subj '/0=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
```

#### 创建 k8s Secret

```shell
kubectl -n istio-system create secret tls cncamp-credential --key=tracing/cncamp.io.key --cert=trancing/cncamp.io.crt
```

#### Istio Gateway 中以 https 方式暴露服务

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
  name: service0
spec:
  selector:
    istio: ingressgateway
  servers:
    - hosts:
        - 'httpsserver.cncamp.io'
      port:
        name: http-service0
        number: 443
        protocol: HTTPS
      tls:
        mode: SIMPLE
        credentialName: cncamp-credential
```

#### 创建 `VirtualService`、`Gateway`

```shell
kubectl -n tracing apply -f tracing/istio-specs.yaml
```

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0omy1i8qij20u902bmxx.jpg" alt="istio-gw" style="zoom:80%;" />

#### 通过 `https://httpsserver.cncamp.io/http/service0` 访问 `service0`

```shell
curl -v -k --resolve httpsserver.cncamp.io:443:10.104.76.93 https://httpsserver.cncamp.io/http/service0
```

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0omymv8g1j20o90dy41h.jpg" alt="https-server-curl" style="zoom:80%;" />

### Jaeger 追踪

#### k8s 中部署 Jaeger

```shell
kubectl apply -f tracing/jaeger.yaml
```

#### 通过访问服务 100 次采样

```shell
for i in $(seq 1 100)
do
  curl -v -s -k --resolve httpsserver.cncamp.io:10.104.76.93 -o /dev/null "https://httpsserver.cncamp.io/http/service0"
done
```

#### 浏览器获取 Jaeger 数据

##### 浏览器访问的方式

- 方式一：修改 tracing svc 类型为 `NodePort`

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0oobefz3ej20uq04fdhs.jpg" alt="tracing-jaeger" style="zoom:80%;" />

- 方式二：`istioctl dashboard jaeger` 获取访问地址 

##### Jaeger 数据

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0omtqzcgej21gh0n0aer.jpg" alt="tracing-service0" style="zoom:80%;" />

<img src="https://tva1.sinaimg.cn/large/e6c9d24egy1h0omu3wdgaj21h808340w.jpg" alt="tracing-service1" style="zoom:80%;" />
