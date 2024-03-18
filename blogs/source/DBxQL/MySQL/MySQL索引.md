# MySQL 索引

索引是数据库高效获取数据的有序数据结构。在数据之外，数据库系统还维护着满足特定查找算法的数据结构，这些数据结构以某种方式引用或指向数据，这样就可以在这些数据结构上实现高级查询算法，这种数据结构就是索引。
索引是关系数据库的内部实现技术，属于内模式的范畴。

优点：
- 提高数据检索效率，降低数据库的IO成本
- 通过索引列对数据进行排序，降低数据排序的成本，降低CPU的消耗

缺点：
- 索引列也要占用空间
- 索引提高了查询效率，但降低了更新的速度，如insert, update, delete

## 索引结构

|索引结构|描述|
| :--- | :--- |
|B+Tree||
|Hash|底层数据结构是用哈希表实现，只有精确匹配索引列的查询才有效，不支持范围查询|
|R-Tree(空间索引)|空间索引是 MyISAM 引擎的一个特殊索引类型，主要用于地理空间数据类型，通常使用较少|
|Full-Text(全文索引)|是一种通过建立倒排索引，快速匹配文档的方式，类似于 Lucene, Solr, ES|

|索引|InnoDB|MyISAM|Memory|
| :--- | :--- | :--- | :--- |
|B+Tree索引|支持|支持|支持|
|Hash索引|不支持|不支持|支持|
|R-Tree索引|不支持|支持|不支持|
|Full-Text索引|5.6版本后支持|支持|不支持|

### B+Tree

二叉树的缺点：顺序插入时，会形成一个链表，查询性能大大降低；大数据量情况下，层级较深，检索速度慢。可以用红黑树来解决。
红黑树也存在大数据量情况下，层级较深，检索速度慢的问题。
为了解决上述问题，可以使用BTree结构。BTree (多路平衡查找树) 以一棵最大度数（max-degree，指一个节点的子节点个数）为5（5阶）的 btree 为例（每个节点最多存储4个key，5个指针）。
> [演示 BTree Visualization](https://www.cs.usfca.edu/~galles/visualization/BTree.html)
> [演示 B+Tree Visualization](https://www.cs.usfca.edu/~galles/visualization/BPlusTree.html)

B+Tree相对于BTree的区别：
- 所有的数据都会出现在叶子节点
- 叶子节点形成一个单向链表

MySQL 索引数据结构对经典的 B+Tree 进行了优化。在原 B+Tree 的基础上，增加一个指向相邻叶子节点的链表指针，就形成了带有顺序指针的 B+Tree，提高区间访问的性能。

### Hash

哈希索引就是采用一定的hash算法，将键值换算成新的hash值，映射到对应的槽位上，然后存储在hash表中。 
如果两个（或多个）键值，映射到一个相同的槽位上，他们就产生了hash冲突（也称为hash碰撞），可以通过链表来解决。

特点：
- Hash索引只能用于对等比较（=、in），不支持范围查询（betwwn、>、<、...）
- 无法利用索引完成排序操作
- 查询效率高，通常只需要一次检索就可以了，效率通常要高于 B+Tree 索引

存储引擎支持情况：
- Memory
- InnoDB: 具有自适应hash功能，hash索引是存储引擎根据 B+Tree 索引在指定条件下自动构建的

### 一些问题

1. 为什么 InnoDB 存储引擎选择使用 B+Tree 索引结构？
   - 相对于二叉树，层级更少，搜索效率高
   - 对于BTree，无论是叶子节点，还是非叶子节点，都会保存数据，但这样导致一页中存储的键值减少，指针也跟着减少，要同样保存大量数据，只能增加树的高度，导致性能降低。
   - 相对于Hash索引，B+Tree支持范围匹配及排序操作

## 索引分类

||含义|特点|关键字|
| :--- | :--- | :--- | :--- |
|主键索引|针对于表中主键创建的索引|默认自动创建，只能有一个|PRIMARY|
|唯一索引|避免同一个表中某数据列中的值重复|可以有多个|UNIQUE|
|常规索引|快速定位特定数据|可以有多个||
|全文索引|全文索引查找的是文本中的关键词，而不是比较索引中的值|可以有多个|FULLTEXT|

在InnoDB存储引擎中，根据索引的存储形式，又可以分为以下两种：
||含义|特点|
| :--- | :--- | :--- |
|聚簇索引(Clustered Index)|将数据存储与索引放在一块，索引结构的叶子节点保存了行数据|必须有，而且只有一个|
|二级索引(Secondary Index)|将数据与索引分开存储，索引结构的叶子节点关联的是对应的主键|可以存在多个|

![image](https://github.com/ghongli/salt2022/assets/7960635/21ef8c2f-fc2f-45e0-a1b3-18bdcaac6c1c)

聚簇索引选取规则：
- 如果存在主键，`主键索引`就是聚簇索引
- 如果不存在主键，将使用`第一个唯一索引`作为聚簇索引
- 如果表没有主键或没有合适的唯一索引，则InnoDB会自动生成一个 rowid 作为隐藏的聚簇索引

聚簇索引整体是一个B+Tree，非叶子节点存放的是键值，叶子节点存放的是行数据，称之为数据页。这就决定了表中的数据也是聚簇索引中的一部分，数据页间通过一个双向链表来链接。
二级索引(普通索引、辅助索引)，是除聚簇索引外的索引，即非聚簇索引。InnoDB的普通索引叶子节点存储的是主键(聚簇索引)的值，而MyISAM的普通索引存储的是记录指针。

聚簇索引就是逻辑顺序和物理顺序保持一致，所以只能建立一个聚簇索引，但可以建立多个唯一索引等。

### 一些问题

1. 参下语句，指出执行效率高的，并说明原因？

   ```mysql
   -- 备注：id为主键，name字段创建的有索引
   select * from user where id = 10;
   select * from user where name = 'Arm';
   ```

   第一条语句执行效率高，因为第二条语句需要回表查询，相当于两个步骤。

2. InnoDB主键索引的B+Tree高度是多少？

   假设一行数据大小为1k，一页中可以存储16行这样的数据。InnoDB的指针占用6个字节的空间，主键假设为bigint，占用字节数为8，可得公式：n*8+(n+1)*6=16*1024，其中8表示bigint占用的字节数，n表示当前节点存储的key的数量，(n+1)表示指针数量(比key多一个)，算出n约为1170。
   如果树的高度为2，则能存储的数据量大概为：1170 * 16 = 18736；如果树的高度为3，则能存储的数据量大概为：1170 * 1170 * 16 = 21939856；
   另外，如果有大量的数据，需要考虑分表、分库。

## 语法

```text
创建索引：如果不加索引类型参数，则创建的是常规索引
CREATE [ UNIQUE | FULLTEXT ] INDEX index_name ON table_name (index_col_name, ...); 
查看索引
SHOW INDEX FROM table_name;
删除索引
DROP INDEX index_name ON table_name;
```

```mysql
-- name字段的值可能会重复，为该字段创建常规索引
create index idx_user_name on tb_user(name);
-- phone字段的值非空，且唯一，为该字段创建唯一索引
create unique index idx_user_phone on tb_user (phone);
-- 为profession, age, status创建联合索引
create index idx_user_pro_age_stat on tb_user(profession, age, status);
-- 为email建立合适的索引来提升查询效率
create index idx_user_email on tb_user(email);
-- 删除索引
drop index idx_user_email on tb_user;
```

## 使用规则

> 索引文件具有B-Tree的最左前缀匹配特性，如果左边的值未确定，那么无法使用此索引。
> 如果有order by的场景，请注意利用索引的有序性。order by 最后的字段是组合索引的一部分，并且放在索引组合顺序的最后，避免出现file_sort的情况，影响查询性能。
> 利用覆盖索引来进行查询操作，来避免回表操作。
> 建组合索引的时候，区分度最高的在最左边。
> count(distinct col) 计算该列除NULL之外的不重复数量。

### 最左前缀法则

如果索引关联了多列(联合索引)，要遵守最左前缀法则。最左前缀法则指的是查询从索引的最左列开始，并且不跳过索引中的列。如果跳跃某一列，索引将部分失效(后面的字段索引失效)。
联合索引中，出现范围查询(<,>)，范围查询右侧的列索引失效。可以用>=、<=来规避索引失效问题。

#### 索引失效情况

1. 在索引列上进行运算操作：`explain select * from tb_user where substring(phone, 10, 2) = '15';`
2. 字符串类型字段使用时，不加引号: `explain select * from tb_user where phone = 19977909915;`
3. 模糊查询中，如果仅尾部模糊匹配(xx%)，索引不会失效；如果是头部模糊匹配(%xx)及搜索模糊匹配(%xx%)，则索引失效；如：`explain select * from tb_user where profession like '%项目';`
4. 用or分割开的条件，如果or其中一个条件的列没有索引，则涉及的索引都不会被用到。
5. 如果MySQL评估使用索引，比全表更慢，则不使用索引。

**SQL提示**是优化数据库的一个重要手段，简单来说，就是在SQL语句中加入一些人为的提示来达到优化操作的目的。

```mysql
-- user index(xx) 使用索引
explain select * from tb_user use index(idx_user_pro) where profession = '项目';
-- ignore index(xx) 不使用那个索引
explain select * from tb_user ignore index(idx_user_pro) where profession = '项目';
-- force index(xx) 必须使用那个索引
explain select * from tb_user force index(idx_user_pro) where profession = '项目';
```

use 是建议，实际使用那个索引，MySQL还会自己权衡运行速度去更改，force 就是无论如何都强制使用该索引。

#### 覆盖索引及回表查询

尽量使用覆盖索引(查询使用了索引，并且需要返回的列，在该索引中已经全部能找到)，减少 select *，避免回表查询，提高效率。

explain extra 字段的含义：
- using index condition：查找使用了索引，但需要回表查询数据 using where **索引下推**；
- using index：查找使用了索引，但需要的数据，都在索引列中能找到，所以不需要回表查询；

如果在聚簇索引(主键、唯一索引等)中直接能找到对应的行，则直接返回行数据，只需要一次查询，那怕是 select *；
如果在辅助索引中找聚簇索引，如`select id,name from tb_user where name='xxx';`，也只需通过辅助索引(name)查找到对应的id，返回name和name索引对应的id即可，只需要一次查询；
如果是通过辅助索引查找其他字段，则需要**回表查询**，如`select id,name,gender from tb_user where name='xxx';` 首先通过普通索引name查找到对应的id，然后通过聚簇索引id，查询id,name,gender的结果集，过程中需要扫描两次索引B+Tree，性能较扫一遍索引树更低。

所以尽量不要用 select *，容易出现回表查询，降级效率，除非有联合索引包含了所有字段。

**issues**: 一张表，有四个字段：id,username,password,status，由于数据量大，需要对SQL进行优化，最优方案是什么？
```mysql
select id,username,password from tb_user where username='arm';
```
给username,password字段，建立联合索引，则不需要回表查询，直接覆盖索引。

#### 前缀索引

当字段类型为字符串（varchar, text等）时，有时需要索引很长的字符串，这会让索引变得很大，查询时，浪费大量的磁盘IO，影响查询效率，此时可以只用字符串的一部分前缀，建立索引，这样可以大大节约索引空间，从而提高索引效率。

`CREATE INDEX index_name ON table_name (columnn(index_col_name), ...); `

前缀长度：可以根据索引的选择性来决定，而选择性是指不重复的索引值(基数)和数据表的记录总数的比值，索引选择性越高，则查询效率越高。唯一索引的选择性是1，这是最好的索引选择性，性能也是最好的。 

选择性公式：
```mysql
select count(distinct email) / count(*) from tb_user;
select count(distinct substring(email, 1, 5)) / count(*) from tb_user;
```

`show index from table_name;` 里的 sub_part，可以看到截取的长度。

#### 单列索引及联合索引

单列索引：一个索引只包含单个列；联合索引：一个索引包含了多个列；在业务场景中，如果存在多个查询条件，考虑针对于查询字段建立索引时，建议建立联合索引，而非单列索引。
`explain select id, phone, name from tb_user where phone = '19977909915' and name = 'Arm'` 语句只会用到phone索引字段。

## 注意事项

- 多条件联合查询时，MySQL优化器会评估那个字段的索引效率更高，会选择该索引完成本次查询。

## 一些原则

1. 针对大数据量，且查询比较频繁的表建立索引；
2. 针对常作为查询条件、排序、分组操作的字段建立索引；
3. 尽量选择区分度高的列作为索引，尽量建立唯一索引；区分度越高，使用索引的效率越高；
4. 如果是字符串类型的字段，字段长度较长，可以针对字段的特点，建立前缀索引；
5. 尽量使用联合索引，减少单例索引；查询时，联合索引很多时间可以覆盖索引，节省存储空间，避免回表，提高查询效率；
6. 要控制索引的数量，索引并不是多多益善，索引越多，唯护索引结构的代价就越大，会影响增删改的效率；
7. 如果索引表不存储NULL值，在建表时使用NOT NULL约束字段。当优化器知道每列是否包含NULL值时，可以更好地确定那个索引最有效地用于查询；

## 索引下推

```mysql
-- 联合索引 idx_user_pro_age_stat(profession, age, status)
explain select id,status from tb_user where profession like '软件工程师%' and age = 25;
```

### < 5.6

检索联合索引 idx_user_pro_age_stat，查询出所有 profession 包含'软件工程师'的主键id，然后通过聚簇索引：主键，判断出所有符合where子句的数据返回，此过程需要回表。

### >= 5.6

检索联合索引 idx_user_pro_age_stat，查询出所有 profession 包含'软件工程师'且age=25的数据，直接返回结果集，无需回表。可见索引下推，是在非主键索引上的优化，可以有效减少回表的次数，大大提升了查询效率。
