[MongoDB][0]

---

### 术语

SCCC Sync Cluster Connection Config

CSRC Config Server Replica Set

### 基本操作

#### Replica Set

复制集（Replica Sets)，是一个基于主/从复制机制的复制功能，进行同一数据的异步同步，从而使多台机器拥有同一数据的都多个副本，由于有自动故障转移和恢复特性，当主库宕机时不需要用户干预的情况下自动切换到其他备份服务器上做主库，一个集群最多可以支持7个服务器，并且任意节点都可以是主节点。所有的写操作都被分发到主节点,而读操作可以在任何节点上进行,可以实现读写分离,提高负载。

```shell
rs.initiate( {
   _id: "csReplSet",
   members: [ { _id: 0, host: "<host>:<port>" } ]
} )

rs.add( { host: <host:port>, priority: 0, votes: 0 } )
rs.status()

var cfg = rs.conf();
cfg.members[0].priority = 1;
cfg.members[1].priority = 1;
cfg.members[2].priority = 1;
cfg.members[3].priority = 1;
cfg.members[0].votes = 1;
cfg.members[1].votes = 1;
cfg.members[2].votes = 1;
cfg.members[3].votes = 1;
rs.reconfig(cfg);

rs.stepDown()
rs.stepDown(600)

rs.remove("<hostname>:<port>")
```

```shell
# create and populate a new collection
use test
var bulk = db.test_collection.initializeUnorderedBulkOp();
people = ["Marc", "Bill", "George", "Eliot", "Matt", "Trey", "Tracy", "Greg", "Steve", "Kristina", "Katie", "Jeff"];
for(var i=0; i<100; i++){
   user_id = i;
   name = people[Math.floor(Math.random()*people.length)];
   number = Math.floor(Math.random()*10001);
   bulk.insert( { "user_id":user_id, "name":name, "number":number });
}
bulk.execute(); 
```



#### Sharded Cluster

```shell
rs.initiate( {
   _id: "configReplSet",
   version: 1,
   configsvr: true,
   members: [ { _id: 0, host: "<host>:<port>" } ]
} )

mongos --configdb configReplSet/mongodb07.example.net:27019,mongodb08.example.net:27019,mongodb09.example.net:27019  --bind_ip localhost,<hostname(s)|ip address(es)>

# add initial replica set as a shard
mongod --replSet "rs1" --shardsvr --port 27017 --bind_ip localhost,<hostname(s)|ip address(es)>

mongo mongodb3.example.net
rs.initiate( {
   _id : "rs1",
   members: [
       { _id: 0, host: "mongodb3.example.net:27017" },
       { _id: 1, host: "mongodb4.example.net:27017" },
       { _id: 2, host: "mongodb5.example.net:27017" }
   ]
})

mongo mongodb6.example.net:27017/admin

sh.addShard( "rs0/mongodb0.example.net:27017,mongodb1.example.net:27017,mongodb2.example.net:27017" )
```



```shell
# 开启分片功能
sh.enableSharding( "test" )
use test
db.test_collection.createIndex( { number : 1 } )
sh.shardCollection( "test.test_collection", { "number" : 1 } )
# confirm the shard iss balancing
db.stats()
db.printShardingStatus()

use config
db.databases.find()
db.databases.find( { "partitioned": true } )

# list shards
db.adminCommand( { listShards : 1 } )
sh.status()
db.printShardingStatus()

sh.addShard( "mongodb0.example.net:27018" )
sh.addShard( "rs1/mongodb0.example.net:27018" )

# remove chunks from the shard
db.adminCommand( { removeShard: "mongodb0" } )
# move db to another primary shard
db.runCommand( { movePrimary: "fizz", to: "mongodb1" })

# disable the balancer
sh.stopBalancer()
db.adminCommand( { balancerStop: 1 } )

# enable the balancer
sh.setBalancerState(true)
db.adminCommand( { balancerStart: 1 } )

# disable balancing on a collection
sh.disableBalancing("students.grades")

# enable balancing on a collection
sh.disableBalancing("students.grades")

# confirm balancing is enabled or disabled
db.getSiblingDB("config").collections.findOne({_id : "students.grades"}).noBalance;

!sh.getBalancerState() && !sh.isBalancerRunning()

sh.getBalancerState()
use config
while(sh.isBalancerRunning()) {
	print("waiting...");
	sleep(1000);
}
```

#### 幂等 Idempotency

##### Downgrade the feature compatibility

```shell
# disable cluster-to-cluster replication and use write blocking
db.runCommand( { setUserWriteBlockMode: 1, global: false } )

# view featureCompatibilityVersion
db.adminCommand( { getParameter: 1, featureCompatibilityVersion: 1 } )

# enable 6.0 Backwards Incompatible Features
db.adminCommand( { setFeatureCompatibilityVersion: "6.0" } )

# disable 6.0 Backwards Incompatible Features
db.adminCommand( { setFeatureCompatibilityVersion: "5.0" } )

# Set Write Concern Timeout
db.adminCommand( {
   setFeatureCompatibilityVersion: "5.0",
   writeConcern: { wtimeout: 5000 }
} )
```

##### Flush

```shell
# forces to flush all pending writes from the storage layer to disk and locks the entire mongod instance to prevent additional writes until the user releases the lock with a corresponding fsyncUnlock.
# Each fsync lock operation increments the lock count, and fsyncUnlock decrements the lock count.

# use fsync to lock the mongod instance and block write operations for the purpose of capturing backups.
# the lock option ensures that the data files are safe to copy using low-level backup utilities such as cp, scp, or tar. 
# An fsync lock is only possible on individual mongod instances of a sharded cluster, not on the entire cluster.
db.runCommand(
   {
     fsync: 1,
     lock: <Boolean>,
     comment: <any>
   }
)

db.adminCommand( { fsync: 1, lock: true } )

# To unlock the mongod, use db.fsyncUnlock():
db.adminCommand(
   {
     fsyncUnlock: 1,
     comment: <any>
   }
)

db.adminCommand( { fsyncUnlock: 1 } )

db.fsyncUnlock()

# Repeat the db.fsyncUnlock() to reduce the lock count to zero to unlock the instance for writes.

# Check Lock Status
serverIsLocked = function() {
	var co = db.currentOp();
	if (co && co.fsyncLock) {
		return true;
	}
	return false;
}
serverIsLocked();
```



### 备份和恢复

> [`mongodump`](https://www.mongodb.com/docs/database-tools/mongodump/#mongodb-binary-bin.mongodump) and [`mongorestore`](https://www.mongodb.com/docs/database-tools/mongorestore/#mongodb-binary-bin.mongorestore) **cannot** be part of a backup strategy for 4.2+ sharded clusters that have sharded transactions in progress, as backups created with [`mongodump`](https://www.mongodb.com/docs/database-tools/mongodump/#mongodb-binary-bin.mongodump) *do not maintain* the atomicity guarantees of transactions across shards.
>
> To capture a point-in-time backup from a sharded cluster you **must** stop *all* writes to the cluster. On a running production system, you can only capture an *approximation* of point-in-time snapshot.



[0]: https://www.mongodb.com/ "mongodb"
[1]: https://www.mongodb.com/docs/	"mongodb docs"
[2]: https://www.mongodb.com/docs/manual/reference/command/ "reference command"

