GTID(Global Transaction ID) 全局事务 ID

---

> 从 MySQL 5.6.5 开始新增了基于 GTID 的复制方式，保证了每个在主库上提交的事务，在集群中有一个唯一的 ID，强化了数据库的主备一致性，故障恢复以及容错能力。

### 概念

GTID 是全局事务 ID，在主库上生成的与事务绑定的唯一标识，不仅在主库上唯一，在集群内也是唯一的；GTID = server_uuid:transaction_id，实际上是由 UUID+TID 组成的。

### 问题

1. `GTID purged`

   ```mysql
   SELECT @@global.gtid_purged; -- 查看已经 purged 的信息
   SELECT @@global.gtid_executed; -- 查看现在的 gtid 信息
   ```

2. `the master has purged binary logs containing GTIDs that the slave requires.`

   ```mysql
   -- 主库
   show global variables like 'gtid_purged';
   -- 从库
   stop slave;
   reset master; // del master all binlog files, 会造成 slave 无法找到 master 的严重后果，慎用！！
   set @@global.gtid_purged='';
   start slave;
   
   -- 观察slave状态
   show slave status\G;
   ```

   

[0]: https://dev.mysql.com/doc/refman/5.7/en/replication-gtids-lifecycle.html

