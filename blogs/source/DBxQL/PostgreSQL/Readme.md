PostgreSQL

---



##### PG 连接参数 `reWriteBatchedInserts`

- PostgreSQL 跟 ORACLE(批量插入，底层使用了 ps.executeBatch) 一样，默认都是支持 jdbc batch 功能；

- 但为了进一步优化性能，PG 在 9.4.1208 版本后，又提供了参数 `reWriteBatchedInserts`，默认是 false；

- 当参数 `reWriteBatchedInserts` true 时，pgjdbc 会将批量的 `insert into ... values(?, ?)`改写为 `insert into ... values(?, ?), (?, ?)`; 

  - 好处：减少了每个 statement 的开销

    JDBC驱动会重写批量insert转换成多行insert，从而限制数据库的调用次数。

    如果使用正确，reWriteBatchedInserts会提升批量insert性能2-3倍。

  - 坏处：如果某个语句执行失败，整个 batch 都会失败

- reWriteBatchedInserts=true 是从postgresql的jdbc高版本驱动(9.4.1209开始加，但是有bug)才加入的特性，建议升级驱动到42.2.2版本，否则即使你将sql写成 insert into values(),(),(),() 这种形式，一样被转化成单条插入。

[1]: https://jdbc.postgresql.org/	"PostgreSQL JDBC Driver"
[2]: https://jdbc.postgresql.org/documentation/94/connect.html	"JDBC Connecting"
[3]: https://jdbc.postgresql.org/documentation/changelog.html#version_42.2.2 "PostgreSQL JDBC 42.2.2 changelog"
[]: https://vladmihalcea.com/postgresql-multi-row-insert-rewritebatchedinserts-property/ "Multi-row inserts with the PostgreSQL"
[5]: https://www.bilibili.com/read/cv15059677 "对比不同数据库对 JDBC batch 的实现细节"

