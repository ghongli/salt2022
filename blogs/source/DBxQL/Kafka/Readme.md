Kafka

---

> kafka 追求性能，可以结合 group 模式，增加 partition 个数，并批量消费队列中的消息；

### Commands

```shell
# 获取主题授权情况
bin/kafka-acls --authorizer-properties zookeeper.connect=localhost:2181 \
 --list --topic test-topic
 
bin/kafka-topics.sh --zookeeper 127.0.0.1:2181 --list
bin/kafka-topics.sh --bootstrap-server 127.0.0.1:9092 --list
bin/kafka-topics.sh --broker-list 127.0.0.1:9092 --list

bin/kafka-server-start.sh config/server.properties
nohup bin/kafka-server-start.sh config/server.properties &

# create a topic
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
# 指定过期时间：retention.ms, retention.minutes, retention.hours
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test –-config retention.ms=180000
bin/kafka-topics.sh --zookeeper localhost:2181 --list
bin/kafka-topics.sh --zookeeper localhost:2181 --desscribe --topic test

kafka-topics.sh --bootstrap-server localhost:9092 --create --topic employee-salary \
    --partitions 1 --replication-factor 1 \
    --config cleanup.policy=compact \
    --config min.cleanable.dirty.ratio=0.001 \
    --config segment.ms=5000
bin/kafka-topics.sh --bootstrap-server localhost:9092 --desscribe \ 
		--topic employee-salary
kafka-console-consumer.sh --bootstrap-server localhost:9092 \
    --topic employee-salary \
    --from-beginning \
    --property print.key=true \
    --property key.separator=,


# send a message to the topic
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test
# from the earliest offset, bing back the original messages.
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning

bin/kafka-topics.sh --zookeeper localhost:2181 --help

# kafka 将消息存储一段时间，并清理超过保留期的消息！
# 确保消息不会超过一定的时间，也不会超过一定的大小
# 消息保留多长时间，通常使用时间，默认 retention.hours=168，即 7 天时间
# 设置更多的时间，将会占用更多的磁盘空间
# log.retention.hours, log.retention.minutes, log.retention.ms 设置了多个，则按最小值进行消息清理
# log.retention.bytes 应用在每个分区上，表示基于保留消息的总字节数，用于消息过期清理，默认 -1 表示不限，使用消息保留时间配置项
grep -i 'log.retention.[hms].*\=' config/server.properties
grep -i 'log.retention.[bytes]*\=' config/server.properties

# alter the topic retention
bin/kafka-configs –alter –zookeeper localhost:2181 –entity-type topics –entity-name test -add-config retention.ms=1800000
bin/kafka-configs --bootstrap-server localhost:9092 –-alter –-entity-type topics –-entity-name test --add-config retention.ms=1800000,retention.bytes=524288000

grep -i 'log.segment.[hms].*\=' config/server.properties
grep -i 'log.segment.[bytes]*\=' config/server.properties
log.segment.bytes 一个存储分段，最大的字节数，默认 1GB
log.segmetn.ms    如果存储分段没有满，等待关闭的时间，默认是一周


log.cleanup.policy=delete|compact
```



### Issues

#### Invalid url in bootstrap.servers

```shell
./kafka-topics.sh --list --bootstrap-server xxx
```

查看 `config/server.properties`  中的 `advertised.listeners` ，如果配置的不是 IP，而是映射值，需要在 `/etc/hosts` 文件中配置映射信息！

配置项：`bootstrap.servers`

```dockerfile
environment:
	KAFKA_ADVERTISED_HOST_NAME: kafka
  KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

确保 `bootstrap.servers` 值在 `config/properties/consumer.properties` or `config/properties/producer.properties`中，设置了正确的 IP、port。

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
# java.net.InetAddress.getCanonicalHostName() if
not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=PLAINTEXT://:9092
advertised.listeners=PLAINTEXT://外网ip:9092
```



[0]: https://kafka.apache.org/27/documentation.html "Kafka 2.7 Documentation"
[1]: https://rmoff.net/2018/08/02/kafka-listeners-explained/	"Kafka Listeners - Explained"

[2]: https://digitalis.io/blog/kafka/kafka-topic-time-based-retention-policies/ "confirming kafka topic time based retention policies"
[3]: https://www.conduktor.io/kafka/kafka-topic-configuration-log-retention "kafka topic configuration: log retention"
[4]: https://www.conduktor.io/kafka/kafka-topics-internals-segments-and-indexes "kafka topic internals: segments and indexes"
