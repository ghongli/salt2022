SQL Server 开启 CDC 详细过程

---

### 示例数据

```sql
-- 执行批处理
SELECT DB_NAME();  
SELECT USER_NAME();  
GO 2 

-- create database
create database example;
go

-- create schema
use example;
go

create schema insurance;
go

-- create table
use example;
go

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE example.insurance.test (
 -- [EntityID] [int] IDENTITY(2,2) NOT NULL,
 [EntityID] [int] IDENTITY(1,2) NOT NULL,
 [name] [nvarchar](10) NULL,
 [type] [float] NULL,
 [address] [varchar](50) NULL,
 [sex] [bit] NULL,
 -- CONSTRAINT [PK_test_even] PRIMARY KEY CLUSTERED
 CONSTRAINT [PK_test_old] PRIMARY KEY CLUSTERED 
(
 [EntityID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

-- 插入示例数据
USE example;
GO

INSERT INTO example.insurance.test
           ([name]
           ,[type]
           ,[address]
           ,[sex])
     VALUES
           ('test'
           ,0.0
           ,'DYT'
           ,'1')
GO 1

-- 查看 cdc 实际信息
-- role_name 列为空(不是 NULL)，表示不限制 role，若有值，需对照授权表 sp_helprolemember 查看用户是否有开启表级 CDC role
exec sys.sp_cdc_help_change_data_capture;

-- 查看具体的授权情况，default role: db_datareader
sp_helprolemember;

-- 授权 NULL TO ReadOnlyUser
exec sp_addrolemember 'NULL','ReadOnlyUser';

-- 回收角色
exec sp_droprolemember 'db_datareader','ReadOnlyUser';
```

### SQL Server 2016 开启 CDC

#### 开启数据库代理服务

##### Windows

1. 服务管理器，`SQL Server 代理(MSSSQLSERVER)`，修改启动类型为自动，并启动服务

##### Linux

```shell
// 查找 mssql-conf 工具
find / -name mssql-conf
// 开启代理服务
mssql-conf set sqlagent.enabled true
```

#### 启用数据库增量复制

```sql
use example
go

EXEC sys.sp_cdc_enable_db
go

-- 检查数据库是否启用增量复制，is_cdc_enabled 1，表示开启成功
SELECT [name], database_id, is_cdc_enabled
FROM sys.databases
WHERE [name] = 'example'
go
```

- `example - table - sys` 中是否存在 `cdc.captured_columns`, `cdc.change_tables`, `cdc.ddl_history`, `cdc.index_columnss`, `cdc.lsn_time_mapping` 表。

#### 表开启增量复制

```sql
use example
go

EXEC sys.sp_cdc_enable_table
@source_schema = N'insurance',
@source_name = N'test',
@role_name = NULL
go

-- 检查表是否启用增量复制，is_tracked_by_cdc 1，表示开启成功
SELECT [name],is_tracked_by_cdc
FROM sys.tables
WHERE [name] IN ('test')
go
```

#### 检查代理是否正确开启 CDC

1. 在 Microsoft SQL Server Management Studio 中查看 "SQL Server 代理 -> 作业"，里面应存在"cdc.example_capture"与"cdc.example_cleanup"，一个负责捕获变化，一个负责清除变化。
2. 在 Microsoft SQL Server Management Studio 中查看 “example -> 表 -> 系统表”中是否存在 "cdc.insurance_test_CT" 表

#### 插入数据验证

```sql
-- 插入示例数据
USE example;
GO

INSERT INTO example.insurance.test
           ([name]
           ,[type]
           ,[address]
           ,[sex])
     VALUES
           ('test'
           ,0.0
           ,'DYT'
           ,'1')
GO 1
```

```sql
-- table data count
SELECT count(1)
FROM [example].[insurance].[test]

-- cdc table data count
SELECT count(1)
FROM [example].[cdc].[insurance_test_CT]
```

### 注：执行 DDL 操作

开启CDC的表，对字段进行了增、删、改等DDL操作时，必须进行如下操作，否则捕获不到 CDC 信息。

1. 关闭表的 CDC

   ```sql
   -- capture_instance一般为schema_table的格式拼接而成，可以通过以下命令，查询实际的值
   exec sys.sp_cdc_help_change_data_capture
   @source_schema = N'insurance',
   @source_name = N'test';
   ```

   ```sql
   -- 禁用 CDC
   use example
   
   go
   EXEC sys.sp_cdc_disable_table
   @source_schema = N'insurance',
   @source_name = N'test',
   @capture_instance = N'insurance_test'
   go
   ```

2. 执行 DDL

   ```sql
   ALTER TABLE example.insurance.test ALTER COLUMN address VARCHAR(130)
   ```

3. 重新开启 CDC

   ```sql
   use example
   go
   
   EXEC sys.sp_cdc_enable_table
   @source_schema = N'insurance',
   @source_name = N'test',
   @role_name = NULL
   go
   ```

### 全库开启 CDC

```sql
-- 全局替换 将 example 替换为实际的数据库名
-- 全局替换 将 INSURANCE 替换为实际的 schema 名称
USE example
GO
EXEC sys.sp_cdc_enable_db
GO

declare @table_name varchar(100)
declare @database_name varchar(100)
declare @schema_name varchar(100)

set @database_name = 'example'
set @schema_name = 'INSURANCE'

declare my_cursor cursor for SELECT TABLE_NAME
                             FROM example.INFORMATION_SCHEMA.TABLES
                             where TABLE_CATALOG = @database_name
                               and TABLE_SCHEMA = @schema_name;
open my_cursor
fetch next from my_cursor into @table_name
while @@FETCH_STATUS = 0
    begin
        begin try
            exec sys.sp_cdc_enable_table
                 @source_schema = @schema_name,
                 @source_name = @table_name,
                 @role_name = NULL
        end try
        begin catch
            print('[ERROR] ' + @table_name)
        end catch

        fetch next from my_cursor into @table_name
    end
close my_cursor
deallocate my_cursor
```

### 全库关闭 CDC

```sql
-- 全局替换 将 example 替换为实际的数据库名
-- 全局替换 将 INSURANCE 替换为实际的 schema 名称
USE example
GO

declare @table_name varchar(100)
declare @database_name varchar(100)
declare @schema_name varchar(100)

set @database_name = 'example'
set @schema_name = 'INSURANCE'

declare my_cursor cursor for SELECT TABLE_NAME
                             FROM example.INFORMATION_SCHEMA.TABLES
                             where TABLE_CATALOG = @database_name
                               and TABLE_SCHEMA = @schema_name;
open my_cursor
fetch next from my_cursor into @table_name
while @@FETCH_STATUS = 0
    begin
        begin try
            EXEC sys.sp_cdc_disable_table
                 @source_schema = @schema_name,
                 @source_name = @table_name,
                 @capture_instance = 'all';
        end try
        begin catch
            print ('[ERROR] ' + @table_name)
        end catch

        fetch next from my_cursor into @table_name
    end
close my_cursor
deallocate my_cursor

EXEC sys.sp_cdc_disable_db
GO
```

## SQLServer AlwaysOn 架构开启CDC

> AlwaysOn, 域，主，从三节点，高可用架构；外部访问指向域节点；

#### 开启库表 CDC 配置

```sql
-- database
exec sys.sp_cdc_enable_db;

-- table
exec sys.sp_cdc_enable_table @source_schema='insurance', @source_name='test', @role_name=NULL;
```

#### 所有节点开启 SQL Server 代理，节点重启时需要开启

#### 添加作业任务

```sql
exec sys.sp_cdc_add_job @job_type = N'capture';
exec sys.sp_cdc_add_job @job_type = N'cleanup';
```

#### 案例

1. 当 AlwaysOn 集群发生故障转移时，辅助节点升级成主节点后，不能捕获到 cdc 事件，需人工重启代理服务后，才可正常捕获 cdc 事件，即对应的 cdc 表才会有数据。

   经过测试验证，主要是由于系统版本为windows 2012R2 ，而数据库版本为 SqlServer 2016 ，系统和数据库版本不一致，兼容性问题导致的。