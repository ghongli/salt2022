Microsoft SQLServer

---

#### Quickstarts

##### Run and Connect - Docker

###### Pull and run the container image

1. Pull the SQL Server Linux container image from Microsoft Container Registry.

   ```bash
   # docker pull mcr.microsoft.com/mssql/server:2019-latest
   docker pull mcr.microsoft.com/mssql/server:2017-latest
   ```

   

2. To run the container image with Docker.

   ```bash
   # SQLServer-2017 Gotapd8-
   docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=<YourStrong@Passw0rd>" \
      -p 1433:1433 --name sql1 --hostname sql1 \
      -v /data/mssql_data:/opt/mssql_data \
      -d mcr.microsoft.com/mssql/server:2017-latest
   ```

   

3. 开启 mssql-server 的代理服务

   ```bash
   docker exec -it sql1 bash
   /opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P "dev@123,"
   /opt/mssql/bin/mssql-conf set sqlagent.enabled true
   exit
   
   docker stop sql1
   docker start sql1
   ```

   

4. 

##### Change the SA password

1. 

Notes:

- GA = General Availability - when the version is officially available and supported.
- CU = Cumulative Update - a periodic release that includes bug fixes, security fixes, and occasionally a small feature.
- CTP = Community Technology Preview - a preview release that precedes the GA of an upcoming new version.
- GDR = General Distribution Release - a release that contains ONLY security fixes.

#### SQL Server 锁的原理

##### 查看会话隔离级别

```sql
DBCC USEROPTIONS
```

默认隔离级别为 `read committed 已提交读`，读数据时会加共享锁。隔离级别如果是已提交读，但需要实现读操作，需要按下面的方式之一解决： 

- 降低隔离级别为未提交读

- 采用不加锁模式读（select 添加 with(nolock), 如select count(1) from a_table with(nolock);）

[0]: https://docs.microsoft.com/zh-cn/sql/t-sql/language-reference?view=sql-server-ver15	"Transact-SQL 引用 - SQL Server"
[1]: https://hub.docker.com/_/microsoft-mssql-server
[2]: https://github.com/Microsoft/mssql-docker
[3]: https://docs.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker "快速入门 - Linux Container"
[4]: https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker "最佳实践 - Docker"

