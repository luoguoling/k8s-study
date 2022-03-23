# java cpu使用率高排查

```bash
#1、根据top命令，发现PID为24379的Java进程占用CPU高达1700%，出现故障。
#2、首先dump出该进程的所有线程及状态
#使用命令 jstack PID 命令打印出CPU占用过高进程的线程栈.
jstack -l 24379 > /tmp/24379.stack
#3、使用top命令找到耗cpu的线程
#使用top -H -p PID 命令查看对应进程是哪个线程占用CPU过高.
top -H -p 24379
#4、将线程的pid 转成16进制，比如28497 = 6f51
$ printf "%x\n" 28497                 
6f51
#5、到第一步dump出来的 24379.stack 里面找6f51 就知道是哪个线程了
grep "6f51" /tmp/24379.stack

```



java常用命令

```bash
jps -v #查看进程详细启动参数
jinfo pid -flags pid  #只输出参数
jstat -gcutil pid 1000 100  #1000毫秒统计一次gc情况，统计100次




```

