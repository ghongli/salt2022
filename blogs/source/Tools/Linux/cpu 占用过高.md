CPU 占用过高

---

### JAVA

```shell
jps 获取 java 进程 pid
jstack pid >> jstack_xx.log 导出CPU占用高进程的线程栈
jstack -l `jps -l | grep xx | awk '{print $1}'` > jstack_xx_`jps -l | grep xx | awk '{print $1}'`_$(date +"%Y%m%d-%H%M%S").log

top -H -p PID 查看对应进程的哪个线程占用CPU过高
echo “obase=16; PID” | bc 将线程的PID转换为16进制,大写转换为小写
从导出的jstack文件中，查找转换成为16进制的线程PID，找到对应的线程栈
分析负载高的线程栈都是什么业务操作，优化程序并处理问题

# 获取线程信息，并找到占用CPU高的线程
ps -mp pid -o THREAD,tid,time | sort -rn 

# 将需要的线程ID转换为16进制格式
printf "%x\n" tid

# 打印线程的堆栈信息
jstack pid |grep tid -A 30

```

