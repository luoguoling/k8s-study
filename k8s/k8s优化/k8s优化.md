# 1.关于dns优化

就是在deployment中去设置dnsPolicy，在不影响集群内服务直接调用的情况下，把ndots从默认的5修改成了2，使代理服务pod在访问server端域名的时候dns解析直接走绝对域名，这样就会避免走 search 域进行匹配，可以正确匹配到ip地址。

https://mp.weixin.qq.com/s/MsKFxfFo8GobhT_OupN1BQ

