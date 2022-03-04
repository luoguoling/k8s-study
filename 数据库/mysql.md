> https://zhuanlan.zhihu.com/p/165641430  mysql主从
>
> https://blog.csdn.net/u013792404/article/details/91591585 mysql主从正式配置
>
> https://blog.csdn.net/u013792404/article/details/94167965  mysql读写分离
>
> https://www.cnblogs.com/jxlwqq/p/5590120.html

> https://www.cnblogs.com/keerya/p/7883766.html  mysql集群

> https://blog.csdn.net/wzt888_/article/details/81639753  mysql mha
>
> https://www.jianshu.com/p/594c27f60200  mycat+keepalived

```sql
#查看错误
show warnning;
#查看主库状态
show master status \G
#查看从库状态
show slave status \G
#重置主记录信息
reset master;
#重置从记录信息
reset slave;
#停止始从
stop slave;
#开始从
start slave;
#清空从所有连接、信息记录
reset slave all;
#删除从
change master to master_host=' ';
 
 
 
 
 
#从库
stop slave;
reset slave all;
show slave status \G
#清除从库配置文件的配置
 
 
 
#主库
reset master;
#清除主库配置文件的配置
#清除mysql.user从库账号
show master status \G
```

