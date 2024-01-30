# MySQL慢查询日志

慢查询日志是 MySQL 提供的一种日志记录，它用来记录在 MySQL 中查询响应时间超过阈值的语句，具体指响应时间超过 long_query_time 值的 SQL，会被记录到慢查询日志。long_query_time 的默认值是 10s，意思是查询响应时间超过 10s 的 SQL语句。

默认情况下，MySQL 是不开启慢查询日志的，需要手动设置这个参数值。当然，如果不是调优需要的话，一般不建议开启慢查询日志，因为开启慢查询日志或多或少会带来一定性能的影响。慢查询日志支持将日志记录写入日志文件，也支持将日志记录写入数据表。

## 慢查询日志参数

slow_query_log：表示是否开启慢查询日志，OFF表示禁用，ON表示开启
slow_query_log_file：MySQL 慢查询日志存储路径。可以不设置该参数，系统会默认给一个缺省值的文件host_name-slow.log
long_query_time：慢查询时间阈值，默认情况下值为 10s
log_queries_not_using_indexes：表示不使用索引的查询超出 long_time_query 的值也会被记录到日志中,默认值是 OFF表示禁用
log_output：表示存储慢查询日志方式，log_output=’FILE’ 表示将日志存入文件，log_output=‘TABLE’ 表示将日志存入数据表 mysql.slow_log

MySQL 同时支持两种日志存储方式，配置的时候以逗号分隔开，如：log_output=’FILE,TABLE’。一般情况下建议将日志记录到文件中，如果将日志记录到数据表中需要耗费更多系统资源。

## 慢查询日志配置

默认情况下，slow_query_log 是禁用的，可以通过设置 slow_query_log 的值开启，如下所示：

```mysql
show variables like '%slow_query_log%';
-- 开启慢查询日志，只对当前数据库生效，如果重启MySQL后，则会失效
-- 如果要永久生效，需要修改配置文件(linux my.cfg, win my.ini)，增加或修改参数slow_query_log=on 和 slow_query_log_file=x/x.log 后，重启MySQL
set global slow_query_log=1;
```

开启了慢查询日志，什么样的SQL才会被记录到日志中？由long_query_time控制，默认情况下 long_query_time 的值为 10s，可以使用命令修改，也可以通过修改配置文件修改。
对于运行时间刚好等于 long_query_time 的情况，是不会被记录下来的，如下：

```mysql
show variables like '%long_query_time%';
set global long_query_time=5;
```

执行修改操作之后，需要重新连接或打开一个会话才能看到修改的值，或者使用 show global variable like ‘%long_query_time%’ 查看。

### 检查对应的慢查询日志

```mysql
-- log_output 用来指定存储日志的方式
show variables like '%log_output%';
-- 设置慢日志存储方式
set global log_output='TABLE';
```

```mysql
use mysql;
set timestamp=xxx;
select sleep(10);
select * from mysql.slow_log;
```

## 调优

1. 系统变量 log_queries_not_using_indexes 未使用索引的查询，也会被记录到慢日志中，如果调优的话，建议开启这个选项，开启这个选项后 index full scan 的 sql 也会被记录到日志中。

```mysql
show variables like '%log_queries_not_using_indexes%';
set global log_queries_not_using_indexes=1;
```

这个开启之后慢查询日志可能会增长的很快，可以设定 log_throttle_queries_not_using_indexes 变量来限制，默认值是 0，也就是不限制，如果该变量值大于 0 如：log_throttle_queries_not_using_indexes = 100，表示每秒记录100条不使用索引的 SQL 语句到慢查询日志中。

2. 默认情况下，管理类的 SQL 语句也不会被记录到慢查询日志中，log_slow_admin_statements 变量表示是否将管理类的 SQL 语句记录到慢查询日中，管理类的 SQL 语句包含：ALTER TABLE, ANALYZE TABLE, CHECK TABLE, CREATE INDEX, DROP INDEX, OPTIMIZE TABLE, and REPAIR TABLE。

3. MySQL 的从库默认不记录慢查询，如果要开启从库的慢查询需要设定 log_slow_slave_statements。

4. 如果要查询有多少条慢查询记录，可以使用系统变量，如下：

```mysql
show global status like '%Slow_queries%';
```


