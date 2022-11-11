kafka UnknownProducerIdException

---

#### 现象

```wiki
org.apache.kafka.common.errors.UnknownProducerIdException

This exception is raised by the broker if it could not locate the producer metadata associated with the producerId in question. This could happen if, for instance, the producer's records were deleted because their retention time had elapsed. Once the last records of the producerId are removed, the producer's metadata is removed from the broker, and future appends by the producer will return this exception.
```

#### 排查



#### 措施

```toml
# Kafka is < 2.4 you can workaround this by increasing the retention time(considering that your system allows that) of your topic's log(e.g 30 days) and to increase the transactional.id.expiration.ms setting( to 24 days)
# log retention 默认 7 天，调整为 30 天
log.retention.hours=720

# transactional.id.expiration.ms 默认为7天，可调整为24天
transactional.id.expiration.ms=2073600000
```



#### 可改进的地方



[1]: https://kafka.apache.org/26/javadoc/org/apache/kafka/common/errors/UnknownProducerIdException.html	"UnknownProducerIdException"
[2]: https://www.e-learn.cn/topic/3683907	"kafka unknown_prooducer_id exception"

