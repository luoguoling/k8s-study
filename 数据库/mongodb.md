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

```

