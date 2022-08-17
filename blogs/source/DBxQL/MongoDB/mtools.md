[mtools][0]
---

### Requisites

#### requirements.txt

```tcl
psutil
# mdb 3.4 use 3.11
pymongo[srv]==3.11
matplotlib
numpy
mtools
```



```sh
pip3 install --upgrade pip
pip3 install -r requirements.txt
```

```sh
cd /data/mtools
curl -O https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-3.4.24.tgz
tar zxf mongodb-linux-x86_64-3.4.24.tgz
ln -fs mongodb-linux-x86_64-3.4.24 mongodb-3.4
```

```sh
vim ~/.bashrc
export MONGODB_3_4_HOME=/data/mtools/mongodb-3.4
export MONGODB_HOME=${MONGODB_3_4_HOME}
export PATH=${MONGODB_HOME}/bin:$PATH
```

```sh
. ~/.bashrc
mongo --version
```

### Init

```sh
mkdir -p /data/mtools/local/mongo-cluster
cd /data/mtools/local/mongo-cluster
```

```sh
mlaunch init --sharded 2 --replicaset --nodes 3 --config 3 --csrs --mongos 3 --port 27050 --noauth
# mlaunch init --binarypath /opt/mongodb-3.4/bin --dir mongodb34/data --replicaset --name dap --nodes 1 --port 27050
# mlaunch init --binarypath /opt/mongodb-3.4/bin --dir mongodb34/data --sharded tic tac --replicaset --name dap --nodes 1 --config 1 --csrs --mongos 3 --port 27050 --auth --auth-roles dbAdminAnyDatabase readWriteAnyDatabase userAdminAnyDatabase clusterAdmin --username admin --password thomas
mlauch list
mlauch stop
mlauch start
```

```sh
mongo --port 27050
db.adminCommand({listShards: 1});
db.adminCommand({buildInfo: 1});
```

```shell
# init
mlaunch init --replicaset --nodes 2 --arbiter --sharded 2 --config 1 --port 19017 --hostname 192.168.1.182

# 启动
mlaunch start

# 停止
mlaunch stop

# mongos 连接
mongodb://192.168.1.182:19017/test

# shard01 连接
mongodb://localhost:19018,localhost:19019,localhost:19020/?replicaSet=shard01

# shard02 连接
mongodb://localhost:19021,localhost:19022,localhost:19023/?replicaSet=shard02
```



### 使用分片集群

```sh
mongo --port 27050
# 开启 database 分片功能
sh.enableSharding('data');
# 对集合执行分片初始化
sh.shardCollection('data.book', {bookId: 'hashed'}, false, {numInitialChunks: 4});
# 查询数据的分布，仅仅 PSS，没有 Arbiter
db.getShardDistribution();
```

```js
db = db.getSiblingDB('data');
var cnt = 0;

for(var i = 0; i < 100; i++) {
    var dl = [];

    for(var j = 0; j < 100; j++) {
        dl.push({
            "bookId": "BBK-"+i+"-"+j,
            "type": 'Revision',
            "version": 'IricSoneVB001',
            "title": "Jackson's Life",
            "subCount": 10,
            "location": "China CN ShenZhen Futian District",
            "author": {
                "name": 50,
                "email": "Foo@yahoo.com",
                "gender": "female",
            },
            "createTime": new Date(),
        });
    }

    cnt += dl.length;
    db.book.insertMany(dl);
    print("insert ", cnt);
}
```



[0]: https://github.com/rueckstiess/mtools/ "github mtools"
[2]: http://blog.rueckstiess.com/mtools/index.html "mtools docs"

