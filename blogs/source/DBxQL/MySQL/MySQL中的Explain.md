# MySQL Explain

## 前言

通常在系统的测试阶段，可以在系统配置文件中开启MySQL中的慢查询日志功能，并且可以设置SQL执行超过多少时间来记录到一个日志文件中。

只要SQL执行的时间超过了设置的时间，就会被记录到日志文件中，可以在日志文件找到执行比较慢的SQL，从而对这些语句进行调优优化，使用 Explain来分析 SQL 语句的性能。

> 默认情况下，MySQL 没有开启慢查询日志，需要手动来设置：[MySQL慢查询日志](./MySQL慢查询日志.md)。
> 如果不需要调优的话，通常不建议启动，因为开启慢查询日志会带来一定的性能影响。

开启慢查询（永久生效）：

```txt
#如果要永久生效，需要修改配置文件；
#[mysqld] 下增加或修改参数，slow_query_log 和 slow_query_log_file，重启即可，如下：

#开启慢查询
slow_query_log=1
#慢查询日志存放的位置
slow_query_log_file=/xxx/mysql-slow.log
#规定慢 SQL 的查询阙值，超过这个值将会被记录到慢查询文件中，单位：秒
long_query_time=3
#慢查询日志以文件或表的形式输出
log_output=FILE
```

## 基本介绍

Explain被称为执行计划，在语句之前增加 explain 关键字，MySQL 会在查询上设置一个标记，模拟MySQL优化器来执行SQL语句，执行查询时，会返回执行计划的信息，并不执行这条SQL。（注意，如果 from 中包含子查询，仍会执行该子查询，将结果放入临时表中）。

- Explain结果是基于数据表中现有数据的。
- Explain结果与MySQL版本有很大的关系，不同版本的优化器的优化策略不同。

语法及示例：

```mysql
-- 语法
explain sql语句；

-- 示例
explain select * from actor;
explain select now() from dual;
```

>  在查询中的每个”表”会输出一行，这里的“表”的意义非常广泛，不仅仅是数据库表，还可以是子查询、一个union 结果等。 

## 列解读

## 概述

|字段|说明|
|;---|;---;|
|id|每个被独立执行的操作标识，标识对象被操作的顺序；id值越大，先被执行；如果相同，可以认为是一组，执行顺序从上到下；如果是子查询，id的序号会递增；|
|table|被操作的对象名称，通常是表名，但有其他格式|
|partitions|匹配的分区信息(对于非分区表值为NULL)|
|type|关联类型或者访问类型，也可以理解成mysql是如何决定查找表中的行，查找数据行的大概范围|
|select_type|连接操作的类型|
|possible_keys|可能用到的索引|
|key|key列显示MySQL实际决定使用的键（索引），必然包含在possible_keys中|
|key_len|被优化器选定的索引键长度，单位是字节|
|ref|表示本行被操作对象的参照对象，无参照对象为NULL|
|rows|查询执行所扫描的元组个数（对于innodb，此值为估计值）|
|extra|执行计划的重要补充信息，当此列出现Using filesort , Using temporary 字样时就要小心了，很可能SQL语句需要优化|

### 几个比较重要的列

#### key 实际决定使用的键

这一列显示MySQL实际采用那个索引来优化对该表的访问。如果没有使用索引，则该列是NULL。如果想强制使用或忽视possible_keys列中的索引，则在查询中使用 force index、ignore index。

#### key_len 选定的索引键长度

这一列显示MySQL在索引里使用的字节数，通过这个值可以算出具体使用了索引中的那些列。

例如，film_actor的联合索引 idx_film_actor_id 由 film_id 和 actor_id 两个int列组成，并且每个int是4字节。通过结果中的key_len=4可推断出查询使用了第一个列：film_id，来执行索引查找。

#### type 关联类型或访问类型，表示查找数据行的大概范围

这一列显示对表的访问方式，表示在表中找到所需行的方式。
常用的类型： ALL、index、range、 ref、eq_ref、const、system、NULL（从左到右，性能从差到好）

- ALL：Full Table Scan， MySQL将遍历全表以找到匹配的行。
- index: Full Index Scan，index与ALL区别为index类型只遍历索引树。
- range: 只检索给定范围的行，使用一个索引来选择行。
- ref: 表示上述表的连接匹配条件，即哪些列或常量被用于查找索引列上的值。
- eq_ref: 类似ref，区别就在使用的索引是唯一索引，对于每个索引键值，表中只有一条记录匹配，简单来说，就是多表连接中使用primary key或者 unique key作为关联条件。
- const、system: 当MySQL对查询某部分进行优化，并转换为一个常量时，使用这些类型访问。如将主键置于where列表中，MySQL就能将该查询转换为一个常量，system是const类型的特例，当查询的表只有一行的情况下，使用system。
- NULL: MySQL在优化过程中分解语句，执行时甚至不用访问表或索引，例如从一个索引列里选取最小值，可以通过单独索引查找完成。

#### extra 补充信息

这一列包含解决查询的详细信息，有以下几种情况：

- Using where: 不用读取表中所有信息，仅通过索引就可以获取所需数据，这发生在对表的全部的请求列都是同一个索引的部分的时候，表示mysql服务器将在存储引擎检索行后再进行过滤。
- Using temporary：表示MySQL需要使用临时表来存储结果集，常见于排序和分组查询，常见 group by ; order by。
- Using filesort：当Query中包含 order by 操作，而且无法利用索引完成的排序操作称为“文件排序”。
- Using join buffer：强调了在获取连接条件时没有使用索引，并且需要连接缓冲区来存储中间结果。如果出现了这个值，那应该注意，根据查询的具体情况可能需要添加索引来改进能。
- Impossible where： 这个值强调了where语句会导致没有符合条件的行（通过收集统计信息不可能存在结果）。
- Select tables optimized away： 这个值意味着仅通过使用索引，优化器可能仅从聚合函数结果中返回一行。
- No tables used： Query语句中使用from dual 或不含任何from子句。

## 总结

- EXPLAIN不会告诉关于触发器、存储过程的信息或用户自定义函数对查询的影响情况
- EXPLAIN不考虑各种Cache
- EXPLAIN不能显示MySQL在执行查询时所作的优化工作
- 部分统计信息是估算的，并非精确值
- EXPALIN只能解释SELECT操作，其他操作要重写为SELECT后查看执行计划
