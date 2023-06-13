Oracle

---

> Oracle 是实例级的；

1. 一些命令

   ```sql
   # 查看当前字符集
   select userenv('language') as characterSet from dual;
   # 修改服务端字符集，需要先关闭数据库，再重新挂载，然后修改字符集，最后重启
   shutdown immediate;
   startup mount;
   alter system enable restricted session;
   alter system set job_queue_processes=0;
   alter system set aq_tm_processes=0;
   alter database open;
   # ora 存在超集、子集的概念，超集不能向子集转变，子集可以向超集转变
   alter database character set INTERNAL_USE ZHS16GBK;
   # 重启，确认字符集
   shutdown immediate;
   startup;
   select userenv('language') as characterSet from dual;
   
   # 查看有关字符集的所有参数
   select * from v$nls_parameters;
   
   # 查看 dump 文件的字符集
   select nls_charset_name(to_number('0354','xxxx')) from dual;
   ```

   

2. 