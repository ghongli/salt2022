MySQL|MariaDB

---

> 活跃线程高一定会带来CPU使用率的增长。MySQL实现中，每一个CPU只能在同一时间内处理一个请求。但是要注意，这里的请求指的是内核层面，而非应用的并发层面。如果排除掉慢查询导致的请求无法正常处理，活跃线程堆积一般都是由于现网业务流量增长造成的。
>
> 在活跃线程达到临界点时，可能在CPU层面开始产生争抢，内核中会产生大量的mutex排他锁，此时性能曲线表现特征为高CPU使用率、高活跃线程、低IO或低QPS。另外一种情况是突然的业务洪峰，建立连接速度非常快，也可能在CPU层面产生大量争抢，从而导致请求堆积。
>
> 在MySQL体系中，内存主要作为数据缓存使用，也就意味着数据需要不断的迭代，最常用是`buffer pool`和`innodb_adaptive_hash_index`内存区域。整个数据库系统的缓存区域，是数据交换最为频繁的位置，如果内存不足和内存页争抢，则会出现各种异常的堆积和慢查询。
>
> 等待锁索引名：DML语句会将锁加在索引行上，所以获取不到的锁一定是在某个索引上。

#### 常用语句

```sql
-- 查看锁表：
select * from information_schema.innodb_trx\G;

-- 查看进程
select * from information_schema.processlist\G;

show binary logs;

show warnings;

show create table tb1;
show index from tb1;

show variables;
show variables like '%char';
show variables like 'character%';

-- 查看SQL语句的执行计划
explain select * from tb1 where name='testname';
-- 如果索引没有被使用，有可能出现了统计信息不准确导致生成了错误的执行计划。
-- 可以通过下述语句重新生成表上的统计信息用以纠正错误计划。
analyze table cdr_db.fin_oper_settle ;

-- 查询引擎状态
show engine innodb status;

set autocommit=0;
set innodb_lock_wait_timeout=1;
start transaction;
-- 查看事务锁信息
select ENGINE_TRANSACTION_ID, index_name, lock_type, lock_mode, LOCK_STATUS, lock_data  from performance_schema.data_locks;

```

#### `mysql_reset_connection()` 的影响

```wiki
mysql_reset_connection();
可以解决长连接而导致的内存占用太大 ，被系统强行杀掉（OOM），从现象看就是MySQL异常重启。
执行后会初始化连接资源，不需要重连和重新做权限验证，但是会将链接恢复到刚刚连接的状态。
造成影响如下：
1. 回滚活跃的事务并重新设置自动提交模式
2. 释放所有表锁
3. 关闭或删除所有的临时表
4. 重新初始化会话的系统变量值
5. 丢失用户定义的设置变量
6. 释放prepare语句
7. 关闭handler变量
8. 将last_insert_id()值设置为0
9. 释放get_lock()获取的锁
10.清空通过mysql_bind_param()调用定义的当前查询属性

```

