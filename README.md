# RPC 框架

- 服务端启动
先保障注册中心启动

```
 cd demo/server
 go run server.go -c config.yaml

```

- 客户端测试

```
cd demo/client
go run client_proxy.go
```

## 设计原理及代码解读
[RPC框架设计单机实现](https://mp.weixin.qq.com/s?__biz=MzIyMzMxNjYwNw==&mid=2247484325&idx=1&sn=5f49b32b1143d97cc1183adbb742607c&chksm=e8215cb5df56d5a3c35b17ee2d5b600492308b95059122d65c129ca5814b80d088344348d1ca&token=1063132055&lang=zh_CN#rd)

[RPC框架设计集群实现](https://mp.weixin.qq.com/s?__biz=MzIyMzMxNjYwNw==&mid=2247484486&idx=1&sn=de03b6d67a4edd4b648c5f9312548c60&chksm=e8215b56df56d240e00d553df7fb624cc678494f74229b0b196f16bfd63938c3c03fd58e4ca9&token=31665640&lang=zh_CN#rd)

扫码关注微信公众号 ***技术岁月*** 支持：

![技术岁月](https://i.loli.net/2021/01/21/orQm9BUkEqKAR6x.jpg)
