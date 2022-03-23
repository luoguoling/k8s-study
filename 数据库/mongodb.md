## mongodb安装

```bash
wget https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-rhel80-5.0.6.tgz
tar zvxf mongodb-linux-x86_64-rhel80-5.0.6.tgz
mv mongodb-linux-x86_64-rhel80-5.0.6 /usr/local/mongodb5
export PATH=/usr/local/mongodb5/bin:$PATH
sudo mkdir -p /var/lib/mongo
sudo mkdir -p /var/log/mongodb
sudo chown `whoami` /var/lib/mongo     # 设置权限
sudo chown `whoami` /var/log/mongodb   # 设置权限
mongod --dbpath /var/lib/mongo --logpath /var/log/mongodb/mongod.log --fork
#或者可以加入开机启动
cd /lib/systemd/system  
vim mongodb.service  
[Unit]
Description=mongodb
After=network.target remote-fs.target nss-lookup.target

[Service]
Type=forking
ExecStart=/usr/local/mongodb5/bin/mongod --config /usr/local/mongodb5/mongodb.conf
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/usr/local/mongodb5/bin/mongod --shutdown --config /usr/local/mongodb5/mongodb.conf
PrivateTmp=true

[Install]
WantedBy=multi-user.target
#修改权限
chmod 754 mongodb.service
#启动服务  
systemctl start mongodb.service  
#关闭服务  
systemctl stop mongodb.service  
#开机启动  
systemctl enable mongodb.service
#开启端口
firewall-cmd --zone=public --add-port=27017/tcp --permanen
firewall-cmd --reload

#mongo常用操作命令
show dbs;
use db;
db.dropDatabase()
db.createCollection("集合名")
show collections/tables
db.集合名.drop()
db.集合名.insert({"键名":键值})
db.集合名.find()
db.集合名.update({"name":"zhangsan"},{$set:{"name":"lisi"}})
db.集合名.remove()
```

## mongodb副本集群搭建

```bash
1）主节点（Primary）
接收所有的写请求，然后把修改同步到所有Secondary。一个Replica Set只能有一个Primary节点，当Primary挂掉后，其他Secondary或者Arbiter节点会重新选举出来一个主节点。默认读请求也是发到Primary节点处理的，需要转发到Secondary需要客户端修改一下连接配置。

（2）副本节点（Secondary）
与主节点保持同样的数据集。当主节点挂掉的时候，参与选主。

（3）仲裁者（Arbiter）
不保有数据，不参与选主，只进行选主投票。使用Arbiter可以减轻数据存储的硬件需求，Arbiter跑起来几乎没什么大的硬件资源需求，但重要的一点是，在生产环境下它和其他数据节点不要部署在同一台机器上。
注意，一个自动failover的Replica Set节点数必须为奇数，目的是选主投票的时候要有一个大多数才能进行选主决策。

#主节点
/usr/local/mongodb5/bin/mongod --port 8319 --dbpath /root/database/mongodb/data/mongo1 --replSet rs1 
#副本节点
/usr/local/mongodb5/bin/mongod --port 8320 --dbpath /root/database/mongodb/data/mongo2 --replSet rs1
#仲裁者
/usr/local/mongodb5/bin/mongod --port 8321 --dbpath /root/database/mongodb/data/mongo3 --replSet rs1

##配置文件可以参考
vi  /usr/local/sdb/mongodb-linux-x86_64-4.0.4/mongodb.conf
port=27017 #端口 
dbpath= /usr/local/sdb/mongodb-linux-x86_64-4.0.4/db #数据库存文件存放目录 
logpath= /usr/local/sdb/mongodb-linux-x86_64-4.0.4/mongodb.log #日志文件存放路径 
logappend=true #使用追加的方式写日志 
fork=true #不以守护程序的方式启用，即不在后台运行
replSet=sciencedb	#Replica Set的名字 集群名称
maxConns=100 #最大同时连接数 
noauth=true #不启用验证 
journal=true #每次写入会记录一条操作日志（通过journal可以重新构造出写入的数据）。 
#即使宕机，启动时wiredtiger会先将数据恢复到最近一次的checkpoint点，然后重放后续的journal日志来恢复。 
storageEngine=wiredTiger #存储引擎有mmapv1、wiretiger、mongorocks 
bind_ip = 10.0.86.193 #这样就可外部访问了，例如从win10中去连虚拟机中的MongoDB


#登录主节点进行配置
mongo --port=8319
use admin
	cfg={ _id:"sciencedb", members:[ {_id:0,host:'127.0.0.1:8319',priority:2}, {_id:1,host:'127.0.0.1:8320',priority:1}, {_id:2,host:'127.0.0.1:8321',arbiterOnly:true}] };
rs.initiate(cfg)	
rs.status() #查看集群状态
#“stateStr” : “PRIMARY”表示主节点, “stateStr” : “SECONDARY”表示从节点， “stateStr” : “ARBITER”,表示仲裁节点
#常用命令

添加节点命令：
添加secondary：rs.add({host: “127.0.0.1:8319”, priority: 1 })
添加仲裁点：rs.addArb(“127.0.0.1:8321”)
移除节点：rs.remove({host: “127.0.0.1:8319”})

```

## 备份恢复

```bash
mongodump -h dbhost -d dbname -o dbdirectory
mongorestore -h <hostname><:port> -d dbname <path>
```

