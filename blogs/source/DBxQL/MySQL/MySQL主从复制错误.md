MySQL 主从复制错误：5.7 server_uuid 相同

>  A slave with the same server_uuid/server_id as this slave has connected to the master;

```sql
A slave with the same server_uuid/server_id as this slave has connected to the master; the first event '' at 4, the last event read from './mysqld-bin.006846' at 15426863, the last byte read from './mysqld-bin.006846' at 15426863. Error code: 1236; SQLSTATE: HY000.
```

- 查看 slave 状态  `show slave status \G;，关注如下的信息：`

  - slave_IO_Running: No
  - Last_IO_Errorno: xx

- 查看 server_id, server_uuid 信息，检查从库是否存在相同的信息

  ```sql
  show variables like '%server%id%';
  show variables like 'server_id';
  show variables like 'server_uuid';
  # 全局唯一的 server_uuid 的一个好处是：可以解决由 server_id 配置冲突带来的 MySQL 主备复制的异常终止
  SHOW GLOBAL VARIABLES LIKE 'server_uuid';
  ```

  如果 server_id 相同，修改 my.cnf 中的 server_id 配置

  如果 server_uuid 相同，原因可能是：

  - 添加新的 slave 时，直接把一个旧的 slave 停掉，然后拷贝到新的机器上启动，结果数据目录的`auth.cnf`里面保存的uuid仍然是旧 slave 的uuid，导致再向master申请binlog时，master神经错乱，无法识别两个slave。

    

    删除 auto.cnf文件(auto.cnf文件在my.cnf中datadir配置的目录下)，然后重启数据库，数据库会重新生成server_uuid和auto.cnf文件。

    

  总结：

  MySQL 5.6 用 128 位的 server_uuid 代替了原本的 32 位 server_id 的大部分功能。原因很简单，server_id 依赖于 my.cnf 的手工配置，有可能产生冲突，而自动产生 128 位 uuid 的算法，可以保证所有的 MySQL uuid 都不会冲突。

  ​	在首次启动时，MySQL 会调用 generate_server_uuid() 自动生成一个 server_uuid，并且保存到 auto.cnf 文件，这个文件目前存在的唯一目的就是保存 server_uuid。

  ​	在 MySQL 再次启动时会读取 auto.cnf 文件，继续使用上次生成的 server_uuid。

  

  在 MySQL 5.6，Slave 向 Master 申请 binlog 时，会首先发送自己的 server_uuid，Master 用 Slave 发送的 server_uuid 代替 server_id (MySQL 5.6 之前的方式)作为 kill_zombie_dump_threads 的参数，终止冲突或者僵死的 BINLOG_DUMP 线程。
  

- server_uuid 相同导致，而且server_uuid是在mysql启动的时候生成

  进入mysql的数据目录，找到 auto.cnf 文件，将此文件移动走，然后重启mysql，生成新的server_uuid 即可。