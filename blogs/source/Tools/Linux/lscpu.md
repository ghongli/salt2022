lscpu

---

### 输出信息

```wiki
Architecture: #架构
CPU op-mode(s): #支持的模式
Byte Order: #字节排序的模式，常用小端模式
CPU(s): #逻辑cpu颗数
On-line CPU(s) list: #在线的cpu数量，有故障或者过热时，某些CPU会停止运行而掉线
Thread(s) per core: #每个核心线程
Core(s) per socket: #每个cpu插槽核数/每颗物理cpu核数
CPU socket(s): #cpu插槽数，即：物理cpu的数量
NUMA node(s): #有几个NUMA节点
Vendor ID: #cpu厂商ID
CPU family: #厂商设定的CPU家族编号
Model: #型号
Model name: #型号名称
Stepping: #步进,可以理解为同一型号cpu的版本号
CPU MHz: #cpu主频
BogoMIPS: #估算MIPS, MIPS是每秒百万条指令
Hypervisor vendor: #虚拟化技术的提供商
Virtualization type: #cpu支持的虚拟化技术的类型 
L1d cache: #一级高速缓存 dcache 用来存储数据
L1i cache: #一级高速缓存 icache 用来存储指令
L2 cache: #二级缓存
L3 cache:  #三级缓存
NUMA node0 CPU(s):   0-3   //四个cpu在同一个numa节点node0上
Flags:  cpu支持的技术特征
```

#### NUMA(Non-Uniform MemoryAccess)

```
中文名：分布式体系结构 （非统一内存体系结构）
与UMA不同，UMA是统一内存体系结构，在UMA中，多个CPU访问内存是没有区别的，成本和速度相同。
而在NUMA中，可以看成每个CPU有自己的内存，被称为本地内存，
CPU之间通过一种方式连结，使得CPU可以访问非管辖范围内的内存（非本地内存）。

因为需要通过另一个通道获取，速度比访问本地内存要慢。
好处是这种方式增加了扩展性。
缺点是速度会受影响，对像mysql这类的数据库软件会有影响。
```

#### 大小端模式

```wiki
Byte Order: Little Endian
小端模式：低位的字节存储在地址较小的位置
大端模式：高位的字节存储在地址较小的位置
判断当前机器的大小端序常用的命令:
lscpu | grep -i byte
```

