使用 Arthas 排查 Java 进程 cpu 飙高的问题

---

>CPU负载过高一般是某个或某几个线程有问题。

1. 快速安装 Arthas

   ```shell
   curl -O https://arthas.aliyun.com/arthas-boot.jar
   java -jar arthas-boot.jar
   
   # 打印帮助信息
   java -jar arthas-boot.jar -h
   
   # 如果下载速度比较慢，使用 aliyun 的镜像
   java -jar arthas-boot.jar --repo-mirror aliyun --use-http
   
   # 卸载
   rm -rf ~/.arthas/
   rm -rf ~/logs/arthas
   ```

   

2. 启动，并选择 jvm 进程

   ```shell
   # 查看详细的进程 COMMAND
   top -Hc -p <pid>
   
   java -jar arthas-boot.jar --repo-mirror aliyun --use-http
   ```

   

3. 筛选线程

   ```shell
   # 查询当前进程中，那个线程占用CPU比较高
   # 显示所有线程的信息，并且把cpu使用率高的线程排在前面
   thread
   
   # 排列出当前进程中，占用CPU比较高前3个线程，相当于 top -Hcp 100 & printf & jstack 三合一的效果
   thread -n 3
   
   # thread id 查看线程堆栈，可以重复执行几次，确保线程执行的方法没有时刻在变化
   thread 108
   
   # 使用 tt 命令监听方法的调用情况，如：
   # 注意：在线上执行这个命令时，一定要记得加上 -n 参数，否则线上巨大的流量可能会瞬间撑爆JVM内存
   tt -t com.google.common.collect.HashBiMap seekByKey -n 100
   
   # 通过 -i 参数后边跟着对应的 INDEX 编号(tt -t 输出的 INDEX)，查看这条记录的详细信息。再通过-w参数，指定一个OGNL表达式，查找相关对象
   tt -i 1000 -w 'target.getApplicationContext()'
   
   # -x 参数指定展开层级，需要将这个参数设置的比环要大一些，才能确保可以发现环路
   tt -i 1000 -w 'target.getApplicationContext()' -x 2
   
   # 查看方法内部的逻辑
   jad com.google.common.collect.HashBiMap seekByKey
   ```

   

4. cpu profiler 火焰图

   ```shell
   # 查看有那些 st 开头的命令
   profiler st
   
   # profiler 默认对 cpu 进行采样
   profiler start
   # 查看目前获取到的样本数量
   profile getSamples
   # 查看状态，运行的时间，具体是指cpu运行的时间
   profiler status
   # 停止采样，并生成火焰图，默认输出格式为 svg，支持 svg, html
   profiler stop
   # 指定输出格式为 html
   profiler stop --format html
   
   # 列出支持的所有事件
   profiler list
   ```

   火焰图是基于 perf 结果产生的 svg 图片，用来展示 cpu 的调用栈。火焰图就是看顶层的那个函数占据的宽度最大。只要有“平顶(plateaus)”，就表示此函数存在性能问题。颜色没有特殊含义，因为火焰图表示的是 cpu 的繁忙程度，所以一般选择暖色调。

   x 轴表示抽样数，如果一个函数在 x 轴占据的宽度超宽，就表示它被抽到的次数多，即执行的时间长。注意：x 轴不代表时间，而是所有的调用栈合并后，按字母顺序排列的。

   y 轴表示调用栈，每一层都是一个函数。调用栈越深，火焰就越高，顶部就是正在执行的函数，下文都是它的父函数。

5. [利用Arthas精准定位Java应用CPU负载过高问题](https://github.com/alibaba/arthas/issues/1202)

   