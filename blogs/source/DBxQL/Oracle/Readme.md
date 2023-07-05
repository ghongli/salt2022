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
   
   # 删除logminer中的输出，仅保留PK与全字段
   alter table ODS_KY.HS_ACCOUNTS_MAIN drop supplemental log group GGS_16290 ;
   alter table ODS_KY.HS_ACCOUNTS_MAIN drop SUPPLEMENTAL LOG DATA (FOREIGN KEY) COLUMNS;
   alter table ODS_KY.HS_ACCOUNTS_MAIN drop SUPPLEMENTAL LOG DATA (UNIQUE) COLUMNS;
   
   # 查看 logminer 的设置
   select supplemental_log_data_min min, supplemental_log_data_pk pk, supplemental_log_data_ui ui, supplemental_log_data_fk fk, supplemental_log_data_all allc from v$database;
   select * from all_log_groups where where TABLE_NAME = 'HS_ACCOUNTS_MAIN';
   ```

   

2. scn、时间转换

   ```sql
   # 时间转换为 scn
   select timestamp_to_scn(to_timestamp('2021-12-21 08:00:00','YYYY-MM-DD HH24:MI:SS')) as scn from dual;
   
   # scn 转换为时间
   select scn_to_timestamp(9389752548) scn from dual;
   ```

   

3. 