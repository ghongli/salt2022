Kafka

---

### Issues

#### kafka-client API 无法生产和消费消息

命令行工具(kafka-console-producer.sh, kafka-console-consumer.sh)能够相互通信，producer 发布的信息，consumer 能够接收。

##### Workround

开启 PLAINTEXT 监听：将 `kafka/config/server.properties` 中的 `advertised.listeners`，设置为 `advertised.listeners=PLAINTEXT://IP:PORT`，重启 kafka。

#### `TimeoutException: Timed out waiting to send the call.`

1. `nc -w 2 -zv ip port` 或 `telnet ip port` 验证网络的连通性

2. 如果 `advertised.listeners` 与 `listeners` 配置属性不同，则发布到 ZooKeeper 供客户端使用的 `listener`

   默认从 java.net.InetAddress.getCanonicalHostName() 返回的值。

   如果未设置，则将使用侦听器的值。与侦听器不同，通告 0.0.0.0 元地址是无效的。 同样与侦听器不同的是，此属性中可以有重复的端口，因此可以将一个侦听器配置为通告另一个侦听器的地址。这在使用外部负载平衡器的某些情况下很有用。

3. `listeners` 由于 Kafka Connect 旨在作为服务运行，因此它还提供了用于管理连接器的 REST API。 可以使用 `listeners` 侦听器配置选项配置 REST API 服务器。

4. `advertised.listeners` 这个配置就是用来给定一个外部地址（外网ip），以便客户端能正确的连接Kafka，这个没有设定配置的话，就会默认使用listeners的值，这样外网的客户端就尝试着去使用内网的ip去连接代理服务器，在外部使用其内网的ip去访问根本无法连的上的，所以一直连不上最后会出现超时异常。

##### Workround

Kafka服务器启动的配置文件 `server.properties` 中正确配置 `advertised.listeners` ，把其配置成外网的ip，这样，外网的客户端就能连接的上Kafka了。

```yaml
############################# Socket Server Settings #############################

# The address the socket server listens on. It will get the value returned from 
# java.net.InetAddress.getCanonicalHostName() if not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=PLAINTEXT://:9092
advertised.listeners=PLAINTEXT://外网ip:9092
```



[0]: https://kafka.apache.org/27/documentation.html "Kafka 2.7 Documentation"
[1]: https://rmoff.net/2018/08/02/kafka-listeners-explained/	"Kafka Listeners - Explained"

