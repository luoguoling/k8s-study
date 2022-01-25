## grafana展示Nginx日志

> https://cloud.tencent.com/developer/article/1802967?from=article.detail.1781705

```bash
#部署服务 Elasticsearch,Logstash,Kibana,Filebeat,grafana
```



### 一.	nginx配置

```bash
# filebeat收集日志格式
     log_format filebeat_nginxlog escape=json '{"@timestamp":"$time_iso8601",'
                                   '"host":"$hostname",'
                                   '"server_ip":"$server_addr",'
                                   #'"client_ip":"$real",'
                                   '"xff":"$http_x_forwarded_for",'
                                   '"domain":"$host",'
                                   '"url":"$request_uri",'
                                   '"referer":"$http_referer",'
                                   '"args":"$args",'
                                   '"upstreamtime":"$upstream_response_time",'
                                   '"responsetime":"$request_time",'
                                   '"request_method":"$request_method",'
                                   '"status":"$status",'
                                   '"size":"$body_bytes_sent",'
                                   '"request_body":"$request_body",'
                                   '"request_length":"$request_length",'
                                   '"protocol":"$server_protocol",'
                                   '"upstreamhost":"$upstream_addr",'
                                   '"file_dir":"$request_filename",'
                                   '"http_user_agent":"$http_user_agent"'
                                   '}';


# server层配置调用log_format,注意配置的自定义变量$real(获取客户端真实IP)

    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
#    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    set $real $http_x_forwarded_for;
    if ( $real ~ (\d+)\.(\d+)\.(\d+)\.(\d+),?(.*) ){
        set $real $1.$2.$3.$4;
    }

    access_log  /data/filebeat_nginxlog/www.api.channel.sdk.sentsss.com/access.log filebeat_nginxlog;

```

### 二.	logstash配置

```bash
#下载GeoLite2-City.mmdb，GeoLite2-City.mmdb是IP信息解析和地理定位的
下载地址：https://github.com/zhengkw/GeoLite2

[root@k8s-master logstash-7.12.1]# cat logstash.conf
input {
    beats {
        host => '0.0.0.0'
        port => 5044 
#	ssl  => false
    }
}

filter {
  #过滤掉CLB健康检查请求日志
  if [http_user_agent] == "clb-healthcheck" {
    drop{}
  }
  geoip {
    #multiLang => "zh-CN"
    target => "geoip"
    source => "client_ip"
    # GeoLite2-City.mmdb文件存放路径/data/logstash-7.12.1/GeoLite2-City.mmdb（用来地图定位和IP解析）
    database => "/data/logstash-7.12.1/GeoLite2-City.mmdb"
    add_field => [ "[geoip][coordinates]", "%{[geoip][longitude]}" ]
    add_field => [ "[geoip][coordinates]", "%{[geoip][latitude]}" ]
    # 去掉显示 geoip 显示的多余信息
    remove_field => ["[geoip][latitude]", "[geoip][longitude]", "[geoip][country_code]", "[geoip][country_code2]", "[geoip][country_code3]", "[geoip][timezone]", "[geoip][continent_code]", "[geoip][region_code]"]
  }
  mutate {
    convert => [ "size", "integer" ]
    convert => [ "status", "integer" ]
    convert => [ "responsetime", "float" ]
    convert => [ "upstreamtime", "float" ]
    convert => [ "[geoip][coordinates]", "float" ]
    # 过滤 filebeat 没用的字段,这里过滤的字段要考虑好输出到es的，否则过滤了就没法做判断
    remove_field => [ "ecs","agent","host","cloud","@version","input","logs_type" ]
  }
  # 根据http_user_agent来自动处理区分用户客户端系统与版本
  useragent {
    source => "http_user_agent"
    target => "ua"
    # 过滤useragent没用的字段
    remove_field => [ "[ua][minor]","[ua][major]","[ua][build]","[ua][patch]","[ua][os_minor]","[ua][os_major]" ]
  }
}
output {
  #stdout { codec => rubydebug}
  elasticsearch {
    hosts => ["http://127.0.0.1:9200"]
    user => "elastic"
    password => "123456"
    index => "logstash-nginx-%{+YYYY.MM}"
  }
}

```

### 三.	filebeat配置

```bash
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /data/filebeat_nginxlog/*/*access.log*
# 日志是json开启这个
  json.keys_under_root: true
  json.overwrite_keys: true
  json.add_error_key: true

output.logstash:
  hosts: ["222.190.107.198:15044"]

processors:
  - add_host_metadata:
      when.not.contains.tags: forwarded
  - add_cloud_metadata: ~
  - add_docker_metadata: ~
  - add_kubernetes_metadata: ~
```

### 四.	grafana配置

```bash
# 安装grafana插件支持展示图表。
grafana-cli plugins install grafana-piechart-panel
grafana-cli plugins install grafana-worldmap-panel

# Grafana 插件地图Worldmap不显示，替换插件里Grafana文件图片地址
cd /var/lib/grafana/plugins/
sed -i 's/https:\/\/cartodb-basemaps{s}.global.ssl.fastly.net\/light_all\/{z}\/{x}\/{y}.png/http:\/\/{s}.basemaps.cartocdn.com\/light_all\/{z}\/{x}\/{y}.png/' \
grafana-worldmap-panel/module.js \
grafana-worldmap-panel/module.js.map
 
sed -i 's/https:\/\/cartodb-basemaps-{s}.global.ssl.fastly.net\/dark_all\/{z}\/{x}\/{y}.png/http:\/\/{s}.basemaps.cartocdn.com\/dark_all\/{z}\/{x}\/{y}.png/'  \
grafana-worldmap-panel/module.js \
grafana-worldmap-panel/module.js.map

# grafana图表模版下载地址：
https://grafana.com/grafana/dashboards/11190

#基本概况
QUERY domain:$domain AND status:$status  查询域名分别对应的pv
```

