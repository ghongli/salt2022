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

   

4. Linux 修改本机时区：`echo 'Asia/Shanghai' >/etc/timezone`

5. Docker SQL Server 时区修改

   1. 重新指定变量启动

      ```bash
      docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=<YourStrong!Passw0rd>' \
         -p 1433:1433 --name sql1 \
         -e 'TZ=Asia/Shanghai'\
         -v /data/files/database/mssql_data:/opt/mssql_data \
         -d mcr.microsoft.com/mssql/server:2017-latest
         
      docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=Gotapd8-' \ 
      		-p 31433:1433 --name SQLServer-2017-31434 \ 
      		-e 'TZ=Asia/Shanghai' \ 
      		-v /data/files/database/mssql_data:/opt/mssql_data \ 
      		-d mcr.microsoft.com/mssql/server:2017-latest
      
      docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=Gotapd8-' -p 31433:1433 --name SQLServer-2017-31434 -e 'TZ=Asia/Shanghai' -v /data/files/database/mssql_data:/opt/mssql_data -d mcr.microsoft.com/mssql/server:2017-latest
      
      docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=Gotapd8-' -p 21433:1433 --name SQLServer-2017-21434 -e 'TZ=UTC' -v /data/files/database/mssql_data:/opt/mssql_data -d mcr.microsoft.com/mssql/server:2017-latest
      ```

      

   2. 进入容器修改

      ```bash
      #进入容器
      docker exec -it  sqlserver bash
      #修改时间
      apt update                    #为了安装tzdata
      apt install tzdata            #为了获取/usr/share/zoneinfo
      rm /etc/localtime
      ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
      date                          #显示为当前时区时间
      
      #注:有些应用从TZ中获取时区,所以还要加上 ENV 
      TZ="Asia/Shanghai"
      # or -----
      TZ=Asia/Shanghai
      ln -fs /usr/share/zoneinfo/$TZ /etc/localtime
      
      #tzselect 命令找到正确的时区值
      tzselect
      #如果选择中国上海  分别输入 4 9 1 1
      
      #最后重启容器 在重新检查日志时间 
      docker restart sqlserver
      docker logs -f sqlserver
      ```

      

6. 

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
   # lsn 信息转时间
   # 根据 lsn 编号在 cdc 表里查询，但查询结果出来是 NULL，说明 LSN 被清理了：
   SELECT sys.fn_cdc_map_lsn_to_time(cast('AAAEaQAANGAABQ==' as xml).value('xs:base64Binary(.)', 'varbinary(max)'));
   // hex 串信息
   SELECT sys.fn_cdc_map_lsn_to_time(cast('00018d58000048990001' as xml).value('xs:hexBinary(.)', 'varbinary(max)'));
   
   # lsn 估计被清理了，使用 EXEC sp_cdc_help_jobs 确认 cleanup retention，结果是 4320分钟(三天)，且没有做过手动清理；
   EXEC sp_cdc_help_jobs;
   
   # 使用 cdc.fn_cdc_get_all_changes_capture_instance ( from_lsn , to_lsn , '<row_filter_option>' ) 备份源表 cdc 表信息(dump 下源表的 cdc 信息，替换 dbo_Transactions 为实际信息)：
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

   

3.  process

   ```sql
   # 从 master.dbo.sysprocesses 查看 process
   select spid, ecid, dbid, status, loginname, hostname, cmd, request_id, cpu from sys.sysprocesses where loginname = 'sa';
   ```

   **Sys.SysProcesses 系统表是一个很重要的系统视图，主要用来定位与解决Sql Server的阻塞和死锁**

   视图中主要的字段：
   1. Spid：Sql Servr 会话ID
   2. Kpid：Windows 线程ID
   3. Blocked：正在阻塞求情的会话 ID。如果此列为 Null，则标识请求未被阻塞
   4. Waittype：当前连接的等待资源编号，标示是否等待资源，0 或 Null表示不需要等待任何资源
   5. Waittime：当前等待时间，单位为毫秒，0 表示没有等待
   6. DBID：当前正由进程使用的数据库ID
   7. UID：执行命令的用户ID
   8. Login_time：客户端进程登录到服务器的时间。
   9. Last_batch：上次执行存储过程或Execute语句的时间。对于系统进程，将存储Sql Server 的启动时间
   10. Open_tran：进程的打开事务个数。如果有嵌套事务，就会大于1
   11. Status：进程ID 状态，dormant = 正在重置回话 ; running = 回话正在运行一个或多个批处理 ; background = 回话正在运行一个后台任务 ; rollback = 会话正在处理事务回滚 ; pending = 回话正在等待工作现成变为可用 ; runnable = 会话中的任务在等待获取 Scheduler 来运行的可执行队列中 ; spinloop = 会话中的任务正在等待自旋锁变为可用 ; suspended = 会话正在等待事件完成
   12. Hostname：建立链接的客户端工作站的名称
   13. Program_name：应用程序的名称，就是 连接字符串中配的 Application Name
   14. Hostprocess：建立连接的应用程序在客户端工作站里的进程ID号
   15. Cmd：当前正在执行的命令
   16. Loginame：登录名

   

   应用实例：

   1. 检查数据库是否发生阻塞

      先查找那个链接的 blocked 字段不为0，如 SPID 53的blocked 字段不为0，而是 52。SPID 52 的 blocked 为0，就可以得出结论：此时有阻塞发生，53 被 52 阻塞住了。

      如果发现一个连接的 blocked 字段的值等于它自己，那说明这个连接正在做磁盘读写，它要等自己的 I/O 做完。

   2. 查找链接在那个数据库上

      检查 dbid 即可。得到 dbid，可以运行以下查询得到数据库的名字：

      select name,dbid from master.sys.sysdatabases;

   3. 查看此进程执行的SQL 是哪个，查找问题原因
      dbcc inputbuffer(spid);

   4. KILL 掉当前导致阻塞的SQL
      kill　spid

   5. sql阻塞进程查询

      ```sql
      select A.SPID as 被阻塞进程,a.CMD AS 正在执行的操作,b.spid AS 阻塞进程号,b.cmd AS 阻塞进程正在执行的操作
      from master..sysprocesses a,master..sysprocesses b
      where a.blocked<>0 and a.blocked= b.spid;
      
      exec sp_who 'active'; --查看系统内所有的活动进程 BLK不为0的为死锁
      
      exec sp_lock 60; --返回某个进程对资源的锁定情况
      
      SELECT object_name(1504685104); --返回对象ID对应的对象名
      
      DBCC INPUTBUFFER (63); --显示从客户端发送到服务器的最后一个语句
      ```

      

   6. 简洁查询正在运行的进程SQL

      ```sql
      SELECT   spid,
               blocked,
               DB_NAME(sp.dbid) AS DBName,
               program_name,
               waitresource,
               lastwaittype,
               sp.loginame,
               sp.hostname,
               a.[Text] AS [TextData],
               SUBSTRING(A.text, sp.stmt_start / 2,
               (CASE WHEN sp.stmt_end = -1 THEN DATALENGTH(A.text) ELSE sp.stmt_end
               END - sp.stmt_start) / 2) AS [current_cmd]
      FROM     sys.sysprocesses AS sp OUTER APPLY sys.dm_exec_sql_text (sp.sql_handle) AS A
      WHERE    spid > 50
      ORDER BY blocked DESC, DB_NAME(sp.dbid) ASC, a.[text];
      ```

      

   7. 查询死锁进程语句

      ```sql
      select
       request_session_id spid, 
       OBJECT_NAME(resource_associated_entity_id) tableName 
       from
       sys.dm_tran_locks 
       where
       resource_type='OBJECT'
      ```

      

   

4.  角色

   - 固定服务器

     | 角色          | 说明                                        |
     | ------------- | ------------------------------------------- |
     | sysadmin      | 执行SQL Server中的任何动作                  |
     | serveradmin   | 配置服务器设置                              |
     | setupadmin    | 安装复制和管理扩展过程                      |
     | securityadmin | 管理登录和CREATE DATABASE的权限以及阅读审计 |
     | processadmin  | 管理SQL Server进程                          |
     | dbcreator     | 创建和修改数据库                            |
     | diskadmin     | 管理磁盘文件                                |

   

   ​		sa登录永远是固定服务器角色syadmin中的成员，并且不能从该角色中删除。

   ​		注意：只有当没有其他方法登录到SQL Server系统中时，再使用sa登录。

   ​		注意：不能添加、修改或删除固定服务器角色。另外，只有固定服务器角色的成员，才能执行sp_addsrvrolemember、sp_dropsrvrolemember，两个系统过程，从角色中添加或删除登录账户。

   ​		使用系统过程sp_srvrolepermission可以浏览每个固定服务器角色的权限。该系统过程的语法形式为：

   ​			sp_srvrolepermission[[@srvrolename =] 'role']

   ​		如果没有指定role的值，那么所有的固定服务器角色的权限都将显示出来。

   

   - 固定数据库

     1. db_owner          可以执行数据库中技术所有动作的用户

     1. db_accessadmin     可以添加、删除用户的用户

     1. db_datareader       可以查看所有数据库中用户表内数据的用户

     1. db_datawriter        可以添加、修改或删除所有数据库中用户表内数据的用户

     1. db_ddladmin        可以在数据库中执行所有DDL操作的用户

     1. db_securityadmin     可以管理数据库中与安全权限有关所有动作的用户

     1. db_backoperator     可以备份数据库的用户(并可以发布DBCC和CHECKPOINT语句，这两个语句一般在备份前都会被执行)

     1. db_denydatareader   不能看到数据库中任何数据的用户

     1. db_denydatawriter    不能改变数据库中任何数据的用户

     在数据库中，每个固定数据库角色都有其特定的权限。这就意味着对于某个数据库来说，固定数据库角色的成员的权限是有限的。使用系统过程sp_dbfixedrolepermission就可以查看每个固定数据库角色的权限。该系统过程的语法为：

     ​	sp_db.xedrolepermission [[@rolename =] 'role']

     如果没有指定role的值，那么所有固定数据库角色的权限都可以显示出来

   - 用户自定义

     

5.  MULTI_USER mode

   ```sql
   #if the database is in Single_User mode
   USE [master];
   GO
   ALTER DATABASE [YourDatabaseNameHere] SET MULTI_USER WITH NO_WAIT;
   GO
   ```

   

6.  

[0]: https://docs.microsoft.com/zh-cn/sql/t-sql/language-reference?view=sql-server-ver15	"Transact-SQL 引用 - SQL Server"
[1]: https://hub.docker.com/_/microsoft-mssql-server
[2]: https://github.com/Microsoft/mssql-docker
[3]: https://docs.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker "快速入门 - Linux Container"
[4]: https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker "最佳实践 - Docker"

[5]: https://docs.microsoft.com/en-us/sql/connect/jdbc/jdbc-driver-support-for-high-availability-disaster-recovery?view=sql-server-ver15 "JDBC driver support for High Availability, disaster recovery"
[6]: https://github.com/microsoft/mssql-jdbc "mssql jdbc"

