SQL Server 自增属性增删原理

---

> 总结：在SQL Server对于现有表的列进行添加或删除自增属性，都是通过临时表作为中转表来实现的；
>
> `set identity_insert db.schema.table ON` 会排他，当一个表开启了，其他的表就会显示插入自增值失败

1. SQL Server 有类似于 alter table 的语法，直接修改表的列为自增列的嘛？

   没有，那么表设计中是如何实现的？

2. 对现有列添加自增属性

   ```sql
   --表若存在，就删除
   if(object_id('t1') is not null)
   begin 
       drop table t1 ;
   end ;
   
   --创建测试表
   create table t1(id int, c1 char(10), c2 char(10)) ; 
   
   --插入测试数据
   insert into t1(id, c1, c2) values(1,'foo','bar'),(10,'foo','bar'),(100,'fo','bar'),(1000,'foo','bar');
   ```

   - 开启SQL Server Profile
   - 打开表设计界面，修改ID列为自增列，保存
   - 停止SQL Server Profile跟踪，查看SQL Server内部实现

   详细步骤如下：

   ```sql
   --1、创建与原表表结构一致的临时表，并且在列上添加了自增属性
   CREATE TABLE dbo.Tmp_t1
       (
       id int NOT NULL IDENTITY (10, 1),
       c1 char(10) NULL,
       c2 char(10) NULL
       )  ON [PRIMARY]
   
   --2、把新增临时表的锁升级为表锁
   ALTER TABLE dbo.Tmp_t1 SET (LOCK_ESCALATION = TABLE)
   
   --3、设置新增临时表的自增列为可插入状态
   SET IDENTITY_INSERT dbo.Tmp_t1 ON
   
   --4、把原表中的数据插入到临时表里
   IF EXISTS(SELECT * FROM dbo.t1)
        EXEC('INSERT INTO dbo.Tmp_t1 (id, c1, c2)
           SELECT id, c1, c2 FROM dbo.t1 WITH (HOLDLOCK TABLOCKX)')
   
   --5、设置新增临时表的自增列为不可插入状态
   SET IDENTITY_INSERT dbo.Tmp_t1 OFF
   
   --6、删除原表
   DROP TABLE dbo.t1
   
   
   --7、把临时表的表名修改为跟原表一致
   EXECUTE sp_rename N'dbo.Tmp_t1', N't1', 'OBJECT' 
   ```

   SQL Server内部也是通过使用临时表作为中转来实现把列修改为自增列的；

   注意：在设计中是修改了标识种子为10，所以在创建临时表Tmp_t1的时候出现了IDENTITY(10，1)，如果没有修改标识种子，默认的是IDENTITY(1，1)，可以在修改完成后使用以下语句进行修改：

   ```sql
   --修改自增列的标识种子
   DBCC CHECKIDENT('t1', reseed, 100) ; 
   
   --查看自增列的当前值
   SELECT IDENT_CURRENT('t1')
   ```

   

3. 对现有列删除自增属性

   - 开启SQL Server Profile
   - 打开表设计界面，修改ID列为非自增列，保存
   - 停止SQL Server Profile跟踪，查看SQL Server内部实现

   详细步骤如下：

   ```sql
   --创建临时表dbo.Tmp_t1
   CREATE TABLE dbo.Tmp_t1
       (
       id int NOT NULL,
       c1 varchar(20) NULL,
       c2 varchar(20) NULL
       )  ON [PRIMARY]
   
   --锁定临时表，锁级别为表锁
   ALTER TABLE dbo.Tmp_t1 SET (LOCK_ESCALATION = TABLE)
   
   --把原来的表的数据插入到临时表dbo.Tmp_t1
   IF EXISTS(SELECT * FROM dbo.t1)
        EXEC('INSERT INTO dbo.Tmp_t1 (id, c1, c2)
           SELECT id, c1, c2 FROM dbo.t1 WITH (HOLDLOCK TABLOCKX)')
   
   --删除原表
   DROP TABLE dbo.t1
   
   --将临时表改名为原表
   EXECUTE sp_rename N'dbo.Tmp_t1', N't1', 'OBJECT' 
   
   --添加索引
   ALTER TABLE dbo.t1 ADD CONSTRAINT
       PK__t1__3213E83F7F60ED59 PRIMARY KEY CLUSTERED 
       (
       id
       ) WITH( STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
   ```

   SQL Server内部也是通过使用临时表作为中转来实现把列修改为非自增列的；

   

   