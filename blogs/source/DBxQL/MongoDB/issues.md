MongoDB Issues

---

1. mdb SCRAM 认证失败

   使用SCRAM，MongoDB验证所提供的用户凭证 [`name`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.user) , [`password`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.credentials) 和 [`authentication database`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.db) 。身份验证数据库是创建用户的数据库，它与用户名一起用于标识用户。

   mdb 连接串中必须包含 username, passowrd, authentication database(e.g. admin)。

2. cannot perform operation: a background operation is currently running for collection

   ```shell
   # 索引分为前台索引和后台索引
   # 杀掉创建索引的进程：找到 lockType: write，在附近会有一下 opid
   db.currentOp() 
   db.killOp(opid)
   ```

   给一个已经存在的大集合的某个字段增加索引的情况。如果集合有很多数据，并且使用了前台索引，此时，创建索引会把这个集合锁起来，所有对这个集合的写入操作都会挂起，直到索引创建完成为止。如果使用的是后台索引，那么创建索引的过程不会影响数据写入。

   因为大集合创建索引有时候可能需要好几个小时，挂起的写入数据会堆积在内存里面，把内存撑爆。

   此时，千万不要重启 MongoDB，因为重启以后，之前没有完成的索引创建操作，依然会自动启动继续创建。

   正确的做法是杀掉创建索引的进程。

3. 

[0]: https://docs.mongoing.com/ "MongoDB-CN-Manual"