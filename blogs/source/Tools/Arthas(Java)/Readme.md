Arthas

---

>Arthas 是Alibaba开源的Java诊断工具，能够帮助我们快速定位线上问题。通过全局视角实时查看应用 load、内存、gc、线程的状态信息，并能在不修改应用代码的情况下，对业务问题进行诊断，包括查看方法调用的出入参、异常，监测方法执行耗时，类加载信息等，大大提升线上问题排查效率。
>
>基本的安装使用可以参考官方文档：https://arthas.aliyun.com/doc

### 特性

1. 后台异步任务

   ```shell
   # 在后台执行命令
   trace Test t >> test.out 2>&1 &
   
   # 查看有那些任务在执行
   jobs
   
   # bg <job-id>或者fg <job-id>可让任务重新开始执行
   # fg命令将任务转到前台
   # bg命令将任务转到后台
   # 将任务暂停
   ctrl + z
   # 停止任务
   ctrl + c
   # linux: 推出终端
   ctrl + d
   # 停止执行
   kill <job-id>
   ```

   

2. 

### 案例

1. [用户案例](https://github.com/alibaba/arthas/issues?q=label%3Auser-case)

### 命令

1. tt(TimeTunnel)

   > 方法执行数据的时空隧道，记录下指定方法每次调用的入参和返回信息、抛出的异常，并能对这些不同的时间下调用进行观测。

   ```shell
   # 使用 tt 命令监听方法的调用情况
   # 注意：在线上执行这个命令时，一定要记得加上 -n 参数，否则线上巨大的流量可能会瞬间撑爆JVM内存
   # -m 参数指定 Class 匹配的最大数量，防止匹配到的 Class 数量太多导致 JVM 挂起，默认值是 50
   tt -t com.google.common.collect.HashBiMap seekByKey -n 100
   
   # 通过 -i 参数后边跟着对应的 INDEX 编号(tt -t 输出的 INDEX)，查看这条记录的详细信息。再通过-w参数，指定一个OGNL表达式，查找相关对象
   tt -i 1000 -w 'target.getApplicationContext()'
   
   # -x 参数指定展开层级，需要将这个参数设置的比环要大一些，才能确保可以发现环路
   tt -i 1000 -w 'target.getApplicationContext()' -x 2
   ```

   

2. watch 查看参数

   ```shell
   watch -h
   
   watch constant.JdbcUtil count params -e -x 1
   watch common.MemoryMessageConsumer valueConverter "{params[0],params[1]}" "params.length==4  && params[1] == 'IV_CSS_MODEL_POSS' " -x 2
   watch connector.db2.cdc.grpc.Db2GrpcLogMiner getCurrentScnStr "{params,throwExp}" -e -x 2
   # 查看 connector.db2.cdc.grpc.Db2GrpcLogMiner#getCurrentScnStr 函数的返回值
   watch connector.db2.cdc.grpc.Db2GrpcLogMiner getCurrentScnStr returnObj -s -x 2
   
   # 查看第一个参数
   watch com.taobao.container.Test test "params[0]"
   Press Ctrl+C to abort.
   
   # 查看第一个参数的 size
   watch com.taobao.container.Test test "params[0].size()"
   Press Ctrl+C to abort.
   
   # 将结果按 name 属性投影
   watch com.taobao.container.Test test "params[0].{ #this.name }"
   Press Ctrl+C to abort.
   
   # 按条件过滤
   watch com.taobao.container.Test test "params[0].{? #this.name == null }" -x 2
   watch com.taobao.container.Test test "params[0].{? #this.name != null }" -x 2
   ## 一定要用 size()>0 才行！
   watch com.taobao.container.Test test "params[0].{? #this.age > 10 }.size()>0" -x 2
   ### 外面用单引号，里面是双引号
   watch com.taobao.container.Test test "params[0].{? #this.deviceKey =="KPmIDmPKMV"}.size()>0" -x 2
   
   ## 过滤后统计
   watch com.taobao.container.Test test "params[0].{? #this.age > 10 }.size()" -x 2
   
   # 判断是否相等
   watch com.demo.Test test 'params[0]=="xyz"'
   watch com.demo.Test test 'params[0]==123L'
   ```

   

2. thread

   ```shell
   # 查询当前进程中，那个线程占用CPU比较高
   # 显示所有线程的信息，并且把cpu使用率高的线程排在前面
   # 默认采样间隔为 200ms，可以使用 -i 200 指定
   # 为了降低统计自身的开销带来的影响，可以把采样间隔拉长一些，比如 5000 毫秒
   # 线程 CPU 使用率 = 线程增量 CPU 时间 / 采样间隔时间 * 100%
   # 默认按照 CPU 增量时间(deltaTime)降序排列，只显示第一页数据，其中 TIME 为线程运行总 cpu 时间。
   thread
   
   # 找出当前阻塞其他线程的线程
   # 目前只支持找出 synchronized 关键字阻塞住的线程
   thread -b
   
   # 查看指定状态的线程
   thread --state WAITING
   
   # 排列出当前进程中，占用CPU比较高前3个线程，相当于 top -Hpc & printf & jstack 三合一的效果
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
   ```

   

3. logger

   ```shell
   # 查看所有 logger 信息
   logger
   # 查看指定名字的 logger 信息
   logger -n org.springframework.web
   
   # 更新 logger level
   logger --name ROOT --level debug
   ```

   

4. dashborad 当前系统的实时数据面板

5. profiler

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

   

   ### 其他

   1. jstack

      ```shell
      jps -l
      # jstack
      jstack <fe pid> | grep 'Connector runner'
      jstack <fe pid> | grep [-A20] <jobId 上面命令找到的>
      # 查线程数
      jstack -l 79600 |grep 'java.lang.Thread.State: '|wc -l
      
      jstack -l 79600 > jstack_fe_79600_$(date +"%Y%m%d-%H%M%S").log
      ```

      ```shell
      # for n in {1..5} ; do echo $n; sleep 10s; done
      pname=agent
      for pid in `jps -l | grep ${pname} | awk '{print $1}'`; do echo $pid; done
      
      # 每 10s 取一次 jstack，取5次，并将相关文件打包
      # pname fe 或 tm
      s=1; e=5; pname=fe
      if [[ $pname == 'fe' ]]
      then
          pname=agent
      fi
      jstack_name_key=$pname
      for ((n=$s;n<=$e;n++))
      do
          for pid in `jps -l | grep ${pname} | awk '{print $1}'`
          do
              echo 'n->'$n',pid->'$pid
              echo 'start time: '$(date +"%Y%m%d-%H%M%S-%N")
              jstack -l $pid > jstack_${jstack_name_key}_$pid_$(date +"%Y%m%d-%H%M%S").log
              
              if [[ $n == $e ]]
              then
                  jstackFile=jstack_${jstack_name_key}_${pid}_$(date +"%Y%m%d-%H%M").tar.gz
                  tar zcf $jstackFile jstack_$jstack_name_key_$pid_*
                  echo $pid' jstack tar file: '$jstackFile
              fi
              echo 'end   time: '$(date +"%Y%m%d-%H%M%S-%N")
              echo '----------'
          done
          
          sleep 10s
      done
      
      jstack -l `jps -l | grep tm | awk '{print $1}'` > jstack_tm_`jps -l | grep tm | awk '{print $1}'`_$(date +"%Y%m%d-%H%M%S").log
      jstack -l `jps -l | grep agent | awk '{print $1}'` > jstack_fe_`jps -l | grep agent | awk '{print $1}'`_$(date +"%Y%m%d-%H%M%S").log
      ```

      

   2. jstat

      ```shell
      # 执行命令，跑个60秒再截图(遇到任务 hang 住了)
      pname=fe
      if [[ $pname == 'fe' ]]
      then
          pname=agent
      fi
      
      jstat -gcutil `jps -l | grep ${pname} | awk '{print $1}'` 1000 1000 
      ```

      

   3. jmap

      ```shell
      # 检查jvm内存占用情况，列出首位的两个, 显示占用资源最多的java类
      jmap -histo:live <pid> |grep agent
      
      # 查看java堆内存使用情况
      jmap -heap <pid> 
      
      pname=fe
      if [[ $pname == 'fe' ]]
      then
          pname=agent
      fi
      
      jmap -histo:live  `jps -l | grep ${pname} | awk '{print $1}'` > heap-live_`jps -l |grep ${pname} | awk '{print $1}'`_$(date +"%Y%m%d-%H%M%S").log
      
      jmap -dump:live,format=b,file=heap-`jps -l |grep ${pname} | awk '{print $1}'`_$(date +"%Y%m%d-%H%M%S").hprof `jps -l | grep ${pname} | awk '{print $1}'`  
      
      jmap -dump:live,format=b,file=/opt/3447.heap <pid>
      jmap -dump:format=b,file=/opt/3447.heap <pid>
      jmap -F -dump:format=b,file=/opt/3447.heap <pid>
      ```

      

   4. other

      ```shell
      # 查句柄数
      lsof -p 79600 | wc -l
      
      # 移除多少天内的文件
      cd logs && find ./ -maxdepth 1 -mtime +3 | grep *.log* | xargs -I {} rm -rf {}
      ```

      