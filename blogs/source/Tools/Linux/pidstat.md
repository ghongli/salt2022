pidstat 监控cpu/mem/io

### 安装

#### Centos

```shell
yum install sysstat

# 查看版本
pidstat -V

# 如果已安装，可以查看所属的rpm包
whereis pidstat
rpm -qf /usr/bin/pidstat
```

一般会默认安装，如未安装则可以用上面的yum命令进行安装

### 最常用的参数

```wiki
-u：监控cpu
-r：监控内存
-d：监控硬盘

-w:  显示每个进程的上下文切换情况
-t:   显示选择任务的线程的统计信息

-p:  指定进程id
```



### 常用示例

```shell
# 监控cpu
pidstat -u
# 监控内存
pidstat -r
# 监控硬盘
pidstat -d
# 上下文切换
pidstat -w
```

#### `pidstat -d` io 使用情况输出的信息

```wiki
PID：进程id
kB_rd/s：每秒从磁盘读取的KB
kB_wr/s：每秒写入磁盘KB
kB_ccwr/s：任务取消的写入磁盘的KB。当任务截断脏的pagecache的时候会发生。
任务取消的写入磁盘的kb数。在任务压缩脏页缓存时，可能发生这种情况。在这种情况下，其他任务发起的IO不会处理。

iodelay： 表示 I/O 的延迟（单位是时钟周期），包括等待同步块 I/O 和换入块 I/O 结束的时间
任务的I/O阻塞延迟，以时钟周期为单位。包括同步块I/O和换入块I/O

COMMAND:task的命令名
```

#### `pidstat -w` 上下文切换输出的信息

```wiki
cswch/s：表示每秒自愿上下文切换（voluntary context switches）的次数
ncswch/s：表示每秒非自愿上下文切换（non voluntary context switches）的次数

这两个概念一定要牢牢记住，因为它们意味着不同的性能问题：
所谓自愿上下文切换，是指进程无法获取所需资源，导致的上下文切换。
比如说， I/O、内存等系统资源不足时，就会发生自愿上下文切换。

而非自愿上下文切换，则是指进程由于时间片已到等原因，被系统强制调度，进而发生的上下文切换。
比如说，大量进程都在争抢 CPU 时，就容易发生非自愿上下文切换。
```

```wiki
自愿上下文切换变多了，说明进程都在等待资源，有可能发生了 I/O 等其他问题；
非自愿上下文切换变多了，说明进程都在被强制调度，也就是都在争抢 CPU，说明 CPU 的确成了瓶颈；
```

#### 显示线程的统计信息

```shell
pidstat -t -p 17700
```

```wiki
TGID:主线程的表示
TID:线程id
%usr：进程在用户空间占用cpu的百分比
%system：进程在内核空间占用cpu的百分比
%guest：进程在虚拟机占用cpu的百分比
%CPU：进程占用cpu的百分比
%wait: 进程或者线程等待的cpu使用率；此值过高，表示出现进程或线程争抢CPU的情况
CPU：处理进程的cpu编号
Command：当前进程对应的命令
```

#### 查看 cpu 的使用情况

```shell
pidstat -u -p 17700
```

> 输出内容与 -t相同，但只显示进程的情况

#### 查看指定进程的内存使用情况

```shell
pidstat -r -p 17700
```

```wiki
PID：进程标识符
Minflt/s:任务每秒发生的次要错误，不需要从磁盘中加载页
				 任务造成的小错误的总数。小错误指的是还不需要从磁盘中加载一个内存页。
Majflt/s:任务每秒发生的主要错误，需要从磁盘中加载页
				 任务造成的大错误的总数。大错误指的是需要从磁盘中加载一个内存页。

VSZ：整个任务使用的虚拟内存大小，以kb为单位
RSS：常驻集合大小，非交换区物理内存,使用KB
		 任务使用的没有被交换的物理内存，以kb为单位

%MEM:任务占用的可用物理内存的比例
Command：task命令名
```

#### 对输出数据做排序

```shell
pidstat -u | sort -k 8 -r
```

```wiki
sort: 排序
-k : 指定排序用哪一列，例子中是第8列:%CPU
-r : 倒序
```

