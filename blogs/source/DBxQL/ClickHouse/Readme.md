ClickHouse

---

> [ClickHouse][0], [Docs][1], [DockerHub][2], [ClickHouse Github, [clickhouse-go]][4]
>
> [ClickHouse blog][5] , [clickhouse-presentations meetup...][7]

### What is ClickHouse

ClickHouse® is a column-oriented database management system (DBMS) for online analytical processing of queries (OLAP).

因为它允许在运行时创建表和数据库、加载数据和运行查询，而无需重新配置或重启服务。将不同的列分别进行存储，使其有效的处理分析查询，得到每秒几亿行的吞吐能力。

> 数据压缩 - 磁盘空间、cpu 消耗之间的压缩编解码器；
>
> 即使是在未压缩的情况下，紧凑的存储数据也是非常重要的，因为解压缩的速度主要取决于未压缩数据的大小。

#### 特性

1. 列式数据管理系统

   ClickHouse不单单是一个数据库， 它是一个数据库管理系统。因为它允许在运行时创建表和数据库、加载数据和运行查询，而无需重新配置或重启服务。

   相较于行式存储，列式存储在查询性能上更优。同时列式存储的数据压缩比更高，更加节省存储空间。

2. 数据压缩

   除了在磁盘空间和CPU消耗之间进行不同权衡的高效通用压缩编解码器之外，ClickHouse还提供针对特定类型数据的[专用编解码器](https://clickhouse.com/docs/zh/sql-reference/statements/create#create-query-specialized-codecs)，这使得ClickHouse能够与更小的数据库(如时间序列数据库)竞争并超越它们。

3. 数据的磁盘存储

   ClickHouse被设计用于工作在传统磁盘上的系统，它提供每GB更低的存储成本，但如果可以使用SSD和内存，它也会合理的利用这些资源。

4. 多核心并行处理

   使用服务器上一切可用的资源，从而以最自然的方式`并行处理大型查询`。

   每个节点只访问本地内存和存储，节点信息交互和节点本身是并行处理的。查询性能好，易于扩展。

5. 多服务器分布式处理

   数据可以保存在不同的shard上，每一个shard都由一组用于容错的replica组成，查询可以并行地在所有shard上进行处理。

6. 支持 SQL

   支持一种[基于SQL的声明式查询语言](https://clickhouse.com/docs/zh/sql-reference/)，在许多情况下与[ANSI SQL标准](https://clickhouse.com/docs/zh/sql-reference/ansi)相同；支持[GROUP BY](https://clickhouse.com/docs/zh/sql-reference/statements/select/group-by), [ORDER BY](https://clickhouse.com/docs/zh/sql-reference/statements/select/order-by), [FROM](https://clickhouse.com/docs/zh/sql-reference/statements/select/from), [JOIN](https://clickhouse.com/docs/zh/sql-reference/statements/select/join), [IN](https://clickhouse.com/docs/zh/sql-reference/operators/in)以及非相关子查询。

7. 向量引擎

   为了高效的使用CPU，数据不仅仅按列存储，同时还按向量(列的一部分)进行处理，这样可以更加高效地使用 CPU。

8. 实时的数据更新

   为了使查询能够快速在主键中进行范围查找，数据总是以增量的方式有序的存储在MergeTree中。因此，数据可以持续不断地高效的写入到表中，并且写入的过程中不会存在任何加锁的行为。

   支持在表中定义主键。近实时数据更新，支持近实时的数据插入、指标聚合以及索引创建。

9. 索引

   按照主键对数据进行排序，可以在几十毫秒以内完成对数据特定值或范围的查找。

10. 适合在线查询

    在没有对数据做任何预处理的情况下，以极低的延迟处理查询，并将结果加载到用户的页面中。

11. 支持近似计算

    允许牺牲数据精度的情况下，对查询进行加速的方法：

    - 用于近似计算的各类聚合函数(AVG, SUM, etc)
    - 基于数据的部分样本进行近似查询，仅会从磁盘检索少部分比例的数据
    - 不使用全部的聚合条件，通过随机选择有限个数据聚合条件进行聚合

12. 支持数据复制和数据完整性

    使用异步的多主复制技术，当数据被写入任何一个可用副本后，系统会在后台将数据分发给其他副本，以保证系统在不同副本上保持相同的数据。在大多数情况下，能在故障后自动恢复，在一些少数的复杂情况下需要手动恢复。

13. 角色的访问控制

    使用SQL查询实现用户帐户管理，并允许[角色的访问控制](https://clickhouse.com/docs/zh/operations/access-rights)，类似于ANSI SQL标准和流行的关系数据库管理系统。

14. Adaptive Join Algorithm

    支持自定义[JOIN](https://clickhouse.com/docs/zh/sql-reference/statements/select/join)多个表，它更倾向于散列连接算法。

15. 局限性

    - 没有完整的事务支持，且更新删除非常慢。
      - All transaction data should be kept in Postgres (or another OLTP database) and ClickHouse should be used for what it does best: OLAP queries. However, we are excited about the [transaction support](https://github.com/ClickHouse/ClickHouse/issues/22086) experiments released in [22.4](https://clickhouse.com/docs/en/whats-new/changelog/) and look forward to experimenting.
    - 缺少高频率、低延迟的修改或删除已有数据的能力，仅能用于批量删除或修改数据，但符合 [GDPR](https://gdpr-info.eu/)。
    - 稀疏索引不适合通过其键，检索单行的查询。

#### 性能

1. 单个大查询的吞吐量

   吞吐量可以使用每秒处理的行数或每秒处理的字节数来衡量。

   对于分布式处理，处理速度几乎是线性扩展的，但这受限于聚合或排序的结果不是那么大的情况下。

   如果数据被放置在page cache中，则一个不太复杂的查询在单个服务器上大约能够以2-10GB／s（未压缩）的速度进行处理（对于简单的查询，速度可以达到30GB／s）。如果数据没有在page cache中的话，那么速度将取决于你的磁盘系统和数据的压缩率。

2. 处理短查询的延迟时间

   如果当前使用的是HDD，在数据没有加载的情况下，查询所需要的延迟可以通过以下公式计算得知： 查找时间（10 ms） * 查询的列的数量 * 查询的数据块的数量。

3. 处理大量短查询的吞吐量

   在相同的情况下，可以在单个服务器上每秒处理数百个查询（在最佳的情况下最多可以处理数千个）。但是由于这不适用于分析型场景。因此建议每秒最多查询100次。

4. 数据的写入性能

   为了提高写入性能，可以使用多个INSERT进行并行写入，这将带来线性的性能提升。建议每次写入不少于1000行的批量写入，或每秒不超过一个写入请求。

#### 应用场景

1. 用户行为分析

   行为分析系统的表可以制作成一张大的宽表，每个表包含大量的列，可以超过一千列。JOIN的形式相对少一点，可以实现路径分析、漏斗分析和路径转化等功能。

2. 流量和监控

   可以将系统和应用监控指标通过流式计算引擎Flink或Spark streaming将监控数据清洗处理以后，实时写入ClickHouse，然后结合Grafana进行可视化展示。

3. 用户画像

   可以将各种用户特征进行数据加工，制作成包含全部用户的一张或多张用户特征表，提供灵活的用户画像分析、支撑广告和圈人等业务需求。

   随着数据时代的发展，各行各业数据平台的体量越来越大，用户个性化运营的诉求也越来越突出，用户标签系统，做为个性化千人千面运营的基础服务，应运而生。如今，几乎所有行业（如互联网、游戏、教育等）都有实时精准营销的需求。通过系统生成用户画像，在营销时通过条件组合筛选用户，快速提取目标群体。

4. 实时 BI 报表

   根据业务需求，可以实时制作一些及时产出的查询灵活的BI报表，实现秒级查询，绝大多数查询能够实时反馈。BI报表包括订单分析、营销效果分析和大促活动分析。

   - 海量数据实时多维查询

     在数亿至数百亿记录规模大宽表，数百以上维度自由查询，响应时间通常在100毫秒以内。让业务人员能持续探索式查询分析，无需中断分析思路，便于深挖业务价值，具有非常好的查询体验。

   - **实时运营监控报表**

     实时分析订单、收入、用户数等核心业务指标；构建用户来源分析系统，跟踪各渠道PV、UV来源。

[0]: https://clickhouse.com/
[1]: https://clickhouse.com/docs/zh "中文文档"
[2]: https://hub.docker.com/r/clickhouse/clickhouse-server/ "DockerHub"
[3]: https://github.com/ClickHouse/ClickHouse "ck github"
[4]: https://github.com/ClickHouse/clickhouse-go "clickhouse-go"
[5]: https://clickhouse.com/blog/ "clickhouse blog"
[6]: https://clickhouse.com/blog/how-quickcheck-uses-clickhouse-to-bring-banking-to-the-unbanked "how quickcheck uses clickhouse to bring banking to the unbanked"
[7]: https://github.com/ClickHouse/clickhouse-presentation "clickhouse-presentations"