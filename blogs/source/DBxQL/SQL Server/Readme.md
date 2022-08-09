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

##### 将用户角色设置为 sysadmin

```sql
-- 以系统管理员身份，将指定用户的角色改为 sysadmin
exec master.dbo.sp_addsrvrolemember @loginame = N'用户名',   @rolename = N'sysadmin'

-- 查看用户角色
select DbRole = g.name, MemberName = u.name, MemberSID = u.sid  
from sys.database_principals u, sys.database_principals g, sys.database_role_members m  
where g.principal_id = m.role_principal_id  
and u.principal_id = m.member_principal_id

-- 查看用户服务器角色
select 用户名 = u.name,管理员权限 = g.name,是否在用 = u.is_disabled,MemberSID = u.sid  
from sys.server_principals u, sys.server_principals g, sys.server_role_members m  
where g.principal_id = m.role_principal_id  
and u.principal_id = m.member_principal_id  
and g.name = 'sysadmin'
```



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

#### Issues

1. [Connection reset issue while connecting to Always On from JDBC](https://github.com/microsoft/mssql-jdbc/issues/1384)

   ```sql
   loginTimeout=30&multisubnetfailover=true
   ```

   

2. cdc lsn 信息查询及备份

   ```sql
   # 根据 lsn 编号在 cdc 表里查询，但查询结果出来是 NULL：
   SELECT sys.fn_cdc_map_lsn_to_time(cast('AAAEaQAANGAABQ==' as xml).value('xs:base64Binary(.)', 'varbinary(max)'));
   
   # lsn 估计被清理了，使用 EXEC sp_cdc_help_jobs 确认 cleanup retention，结果是 4320分钟(三天)，且没有做过手动清理；
   EXEC sp_cdc_help_jobs;
   
   # 使用 cdc.fn_cdc_get_all_changes_capture_instance ( from_lsn , to_lsn , '<row_filter_option>' ) 备份源表 dbo_Transactions cdc 表信息：
   DECLARE @from_lsn binary(10), @to_lsn binary(10);  
   SET @from_lsn = sys.fn_cdc_get_min_lsn('dbo_Transactions');  
   SET @to_lsn   = sys.fn_cdc_get_max_lsn();  
   SELECT * FROM cdc.fn_cdc_get_all_changes_dbo_Transactions (@from_lsn, @to_lsn, N'all');
   
   -- 更改数据保留
   -- EXECUTE sys.sp_cdc_change_job;
   -- 更改数据保留时间
   EXECUTE sys.sp_cdc_change_job @job_type = N'cleanup', @retention=50;
   -- 停用作业
   EXEC sys.sp_cdc_stop_job @job_type = N'cleanup';
   -- 启用作业
   EXEC sys.sp_cdc_start_job @job_type = N'cleanup';
   -- 再次查看
   EXEC sp_cdc_help_jobs;
   ```

   

3.  

[0]: https://docs.microsoft.com/zh-cn/sql/t-sql/language-reference?view=sql-server-ver15	"Transact-SQL 引用 - SQL Server"
[1]: https://hub.docker.com/_/microsoft-mssql-server
[2]: https://github.com/Microsoft/mssql-docker
[3]: https://docs.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker "快速入门 - Linux Container"
[4]: https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker "最佳实践 - Docker"

[5]: https://docs.microsoft.com/en-us/sql/connect/jdbc/jdbc-driver-support-for-high-availability-disaster-recovery?view=sql-server-ver15 "JDBC driver support for High Availability, disaster recovery"
[6]: https://github.com/microsoft/mssql-jdbc "mssql jdbc"

