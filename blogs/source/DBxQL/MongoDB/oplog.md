# MongoDB oplog

## 简介

oplog是local库下的一个固定集合，Secondary就是通过查看Primary 的oplog这个集合来进行复制的。每个节点都有oplog，记录这从主节点复制过来的信息，这样每个成员都可以作为同步源给其他节点。

Oplog 可以说是Mongodb Replication的纽带了。

## 副本集数据同步的过程

副本集中数据同步的详细过程：Primary节点写入数据，Secondary通过读取Primary的oplog得到复制信息，开始复制数据并且将复制信息写入到自己的oplog。如果某个操作失败（只有当同步源的数据损坏或者数据与主节点不一致时才可能发生），则备份节点停止从当前数据源复制数据。如果某个备份节点由于某些原因挂掉了，当重新启动后，就会自动从oplog的最后一个操作开始同步，同步完成后，将信息写入自己的oplog，由于复制操作是先复制数据，复制完成后再写入oplog，有可能相同的操作会同步两份，不过MongoDB在设计之初就考虑到这个问题，将oplog的同一个操作执行多次，与执行一次的效果是一样的。

作用：
1. 当Primary进行写操作的时候，会将这些写操作记录写入Primary的Oplog 中，而后Secondary会将Oplog 复制到本机并应用这些操作，从而实现Replication的功能。
2. 同时由于其记录了Primary上的写操作，故还能将其用作数据恢复。
3. 可以简单的将其视作MySQL中的binlog。

## oplog的增长速度

oplog是固定大小，他只能保存特定数量的操作日志，通常oplog使用空间的增长速度跟系统处理写请求的速度相当，如果主节点上每分钟处理1KB的写入数据，那么oplog每分钟大约也写入1KB数据。如果单次操作影响到了多个文档（比如删除了多个文档或者更新了多个文档）则oplog可能就会有多条操作日志。db.testcoll.remove() 删除了1000000个文档，那么oplog中就会有1000000条操作日志。如果存在大批量的操作，oplog有可能很快就会被写满了。

大小：
1. Oplog 是一个capped collection。
2. 在64位的Linux, Solaris, FreeBSD, and Windows 系统中，Mongodb默认将其大小设置为可用disk空间的5%（默认最小为1G，最大为50G），或也可以在mongodb复制集实例初始化之前将mongo.conf中oplogSize设置为我们需要的值。

local.oplog.rs 一个capped collection集合.可在命令行下使用--oplogSize 选项设置该集合大小尺寸.
但是由于Oplog 其保证了复制的正常进行，以及数据的安全性和容灾能力。

## oplog注意事项：

[local.oplog.rs](http://docs.mongoing.com/manual-zh/reference/local-database.html#local.oplog.rs)特殊的集合。用来记录Primary节点的操作。

为了提高复制的效率，复制集中的所有节点之间会相互的心跳检测（ping）。每个节点都可以从其他节点上获取oplog。

oplog中的一条操作。不管执行多少次效果是一样的

## oplog的大小

第一次启动复制集中的节点时，MongoDB会建立Oplog,会有一个默认的大小，这个大小取决于机器的操作系统

rs.printReplicationInfo()     查看 oplog 的状态，输出信息包括 oplog 日志大小，操作日志记录的起始时间。

db.getReplicationInfo()   可以用来查看oplog的状态、大小、存储的时间范围。

capped collection是MongoDB中一种提供高性能插入、读取和删除操作的固定大小集合，当集合被填满的时候，新的插入的文档会覆盖老的文档。

所以，oplog表使用capped collection是合理的，因为不可能无限制的增长oplog。MongoDB在初始化副本集的时候都会有一个默认的oplog大小：

- 在64位的Linux,Solaris,FreeBSD以及Windows系统上，MongoDB会分配磁盘剩余空间的5%作为oplog的大小，如果这部分小于1GB则分配1GB的空间
- 在64的OS X系统上会分配183MB
- 在32位的系统上则只分配48MB

oplog的大小设置是值得考虑的一个问题，如果oplog size过大，会浪费存储空间；如果oplog size过小，老的oplog记录很快就会被覆盖，那么宕机的节点就很容易出现无法同步数据的现象。

比如，基于上面的例子，停掉一个备份节点，然后通过主节点插入以下记录，然后查看oplog，发现以前的oplog已经被覆盖了。

通过MongoDB shell连接上这个节点，会发现这个节点一直处于RECOVERING状态。

### 解决方法：

数据同步
在副本集中，有两种数据同步方式：

- initial sync（初始化）：这个过程发生在当副本集中创建一个新的数据库或其中某个节点刚从宕机中恢复，或者向副本集中添加新的成员的时候，默认的，副本集中的节点会从离它最近的节点复制oplog来同步数据，这个最近的节点可以是primary也可以是拥有最新oplog副本的secondary节点。
该操作一般会重新初始化备份节点，开销较大
- replication（复制）：在初始化后这个操作会一直持续的进行着,以保持各个secondary节点之间的数据同步。

#### initial sync
当遇到上面例子中无法同步的问题时，只能使用以下两种方式进行initial sync了

- 第一种方式就是停止该节点，然后删除目录中的文件，重新启动该节点。这样，这个节点就会执行initial sync
注意：通过这种方式，sync的时间是根据数据量大小的，如果数据量过大，sync时间就会很长
同时会有很多网络传输，可能会影响其他节点的工作
- 第二种方式，停止该节点，然后删除目录中的文件，找一个比较新的节点，然后把该节点目录中的文件拷贝到要sync的节点目录中
通过上面两种方式中的一种，都可以重新恢复节点，改变一直处于RECOVERING状态的错误。

## oplog数据结构

下面来分析一下oplog中字段的含义，通过下面的命令取出一条oplog：

db.oplog.rs.find().skip(1).limit(1).toArray()
- ts: 8字节的时间戳，由4字节unix timestamp + 4字节自增计数表示。这个值很重要，在选举(如master宕机时)新primary时，会选择ts最大的那个secondary作为新primary
- op：1字节的操作类型
   - "i"： insert
   - "u"： update
   - "d"： delete
   - "c"： db cmd
   - "db"：声明当前数据库 (其中ns 被设置成为=>数据库名称+ '.')
   - "n": no op,即空操作，其会定期执行以确保时效性
- ns：操作所在的namespace
- o：操作所对应的document，即当前操作的内容（比如更新操作时要更新的的字段和值）
- o2: 在执行更新操作时的where条件，仅限于update时才有该属性

通过"db.printReplicationInfo()"命令可以查看oplog的信息，字段说明：

- configured oplog size： oplog文件大小
- log length start to end: oplog日志的启用时间段
- oplog first event time: 第一个事务日志的产生时间
- oplog last event time: 最后一个事务日志的产生时间
- now: 现在的时间

通过"db.printSlaveReplicationInfo()"可以查看slave的同步状态

副本节点中执行db.printSlaveReplicationInfo()命令可以查看同步状态信息

- source——从库的IP及端口
- syncedTo——当前的同步情况，延迟了多久等信息

当我们插入一条新的数据，然后重新检查slave状态时，就会发现sync时间更新了。