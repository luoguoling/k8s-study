https://cloud.tencent.com/developer/article/1781705

# elk+kafka安装

## 一.elastic安装



### 1.1.安装elastic

```bash
#下载代码
cd /usr/local && wget 121.5.106.77/elasticsearch-7.16.3.tar.gz  && tar zvxf elasticsearch-7.16.3.tar.gz
#创建目录
mkdir -p /data/elasticsearch/{data,logs}
```

### 1.2 配置文件

```yaml
cluster.name: es.kcwl.com
node.name: VM-12-9-centos
path.data: /data/elasticsearch/data
path.logs: /data/elasticsearch/logs
bootstrap.memory_lock: true
network.host: 10.0.12.9
http.port: 9200
discovery.seed_hosts: ["10.0.12.9"]
cluster.initial_master_nodes: ["10.0.12.9"]
http.cors.enabled: true
http.cors.allow-origin: "*"
http.cors.allow-headers: Authorization
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
```

### 1.3调整参数并启动

```bash
#设置jvm参数
vi config/jvm.options
# 根据环境设置，-Xms和-Xmx设置为相同的值，推荐设置为机器内存的一半左右
-Xms512m 
-Xmx512m
#创建普通用户
useradd -s /bin/bash -M es
chown -R es.es /opt/elasticsearch-7.10.2
chown -R es.es /data/elasticsearch/
#调整文件描述符
vim /etc/security/limits.d/es.conf
es hard nofile 65536
es soft fsize unlimited
es hard memlock unlimited
es soft memlock unlimited
#调整内核参数
sysctl -w vm.max_map_count=262144
echo "vm.max_map_count=262144" > /etc/sysctl.conf
sysctl -p
#启动es服务
su -c "/usr/local/elasticsearch-7.16.3/bin/elasticsearch -d" es
#查看服务
[root@VM-12-9-centos elasticsearch-7.16.3]# netstat -tunlp|grep 9200
tcp6       0      0 10.0.12.9:9200          :::*                    LISTEN      424891/java   
```

## 二.部署kafka

### 2.1安装zookeeper

```bash
#下载代码
cd /usr/local && wget 121.5.106.77/zookeeper-3.4.14.tar.gz  && tar zvxf zookeeper-3.4.14.tar.gz
mkdir -pv /data/zookeeper/data /data/zookeeper/logs
```

### 2.2配置文件

```bash
cat /opt/zookeeper/conf/zoo.cfg

tickTime=2000
initLimit=10
syncLimit=5
dataDir=/data/zookeeper/data
dataLogDir=/data/zookeeper/logs
clientPort=2181
```

### 2.3启动检查

```bash
[root@VM-12-9-centos zookeeper-3.4.14]# more start.sh 
./bin/zkServer.sh start
[root@VM-12-9-centos zookeeper-3.4.14]# netstat -tunlp|grep 2181
tcp6       0      0 :::2181                 :::*                    LISTEN      3500/java        
```

## 三 安装kafka

### 3.1 安装kafka

```bash
#下载代码
cd /usr/local && wget 121.5.106.77/kafka_2.12-2.2.0.tar.gz  && tar zvxf kafka_2.12-2.2.0.tar.gz
mkdir -pv  /data/kafka/logs
```

### 3.2配置文件

```yaml
broker.id=0
listeners=PLAINTEXT://10.0.12.9:9092
num.network.threads=3
num.io.threads=8
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600
log.dirs=/data/kafka/logs
num.partitions=1
num.recovery.threads.per.data.dir=1
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1
log.flush.interval.messages=10000
log.flush.interval.ms=1000
log.retention.hours=168
log.segment.bytes=1073741824
log.retention.check.interval.ms=300000
zookeeper.connect=localhost:2181
zookeeper.connection.timeout.ms=6000
group.initial.rebalance.delay.ms=0
delete.topic.enable=true
host.name=VM-12-9-centos
```

### 3.3 启动

```bash
root@VM-12-9-centos kafka_2.12-2.2.0]# more start.sh 
bin/kafka-server-start.sh -daemon config/server.properties
[root@VM-12-9-centos kafka_2.12-2.2.0]# netstat -tunlp|grep 9092
tcp6       0      0 10.0.12.9:9092          :::*                    LISTEN      6248/java    
#创建topic
./bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic shengxian
#查看topic
./bin/kafka-topics.sh -zookeeper 127.0.0.1:2181 --list
#测试生产消费
./bin/kafka-console-producer.sh --broker-list 10.0.12.9:9092 --topic testtopic
./bin/kafka-console-consumer.sh --bootstrap-server 10.0.12.9:9092 --topic testtopic --from-beginning
```

## 四.安装logstash

### 4.1安装logstash

```bash
cd /usr/local && wget 121.5.106.77/logstash-7.10.2.tar.gz  && tar zvxf logstash-7.10.2.tar.gz
```

### 4.2配置文件

```bash
input {
  kafka {
    bootstrap_servers => "10.0.12.9:9092"
    topics => ["shengxian"]
    group_id => "shengxian"
    codec => "json"
  }
}
filter {
  json {
    source => "message"
  }
}
output {
  if [filetype] == "web"{
    elasticsearch {
      hosts => ["http://10.0.12.9:9200"]
      user => "elastic"
      password => "lvcheng@2015"
      index => "sx_webapp-%{+YYYY.MM}"
    }
  }
  if [filetype] == "system"{
    elasticsearch {
      hosts => ["http://10.0.12.9:9200"]
      user => "elastic"
      password => "lvcheng@2015"
      index => "system-%{+YYYY.MM}"
    }

  }
}
```

### 4.3启动检查

```bash
nohup ./bin/logstash -f config/logstash-sample.conf &
```

## 5.安装kibana

### 5.1安装kibana

```bash
cd /usr/local && wget 121.5.106.77/kibana-7.16.3-linux-x86_64.tar.gz  && tar zvxf kibana-7.16.3-linux-x86_64.tar.gz
```

### 5.2配置文件

```yaml
server.port: 5601
server.host: "0.0.0.0"
elasticsearch.hosts: ["http://10.0.12.9:9200"]
elasticsearch.username: "elastic"
elasticsearch.password: "lvcheng@2015"
i18n.locale: "zh-CN"
```

### 5.3启动

```bash
nohup ./bin/kibana -c config/kibana.yml --allow-root &
```

## 六.安装filebeat

### 6.1安装filebeat

```bash
cd /usr/local && wget 121.5.106.77/filebeat-7.16.3-linux-x86_64.tar.gz  && tar zvxf filebeat-7.16.3-linux-x86_64.tar.gz
```

### 6.2配置文件

```yaml
filebeat.inputs:
- type: log
  fields_under_root: true
  fields:
    filetype: web
  paths:
    - /var/log/nginx/access.log
  scan_frequency: 120s
  max_bytes: 10485760
  multiline.pattern: ^\d{2}
  multiline.negate: true
  multiline.match: after
  multiline.max_lines: 100

output.kafka:
  hosts: ["10.0.12.9:9092"]
  topic: shengxian
  version: 2.0.0      
  required_acks: 0
  max_message_bytes: 10485760
```

### 6.3启动

```bash
nohup ./filebeat -e -c /usr/local/filebeat-7.16.3-linux-x86_64/filebeat.yml &
```

## 七 .创建es密码

```bash
bin/elasticsearch-setup-passwords interactive
```

## 八.添加域名解析

```bash
server {
    listen       80;
    server_name  kibana.com.cn;

    client_max_body_size 1000m;

    location / {
        proxy_pass http://172.16.90.6:5601;
    }
}
```

