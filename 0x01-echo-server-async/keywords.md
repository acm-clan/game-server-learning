1 tcp server and client
2 log
3 string dump
4 add bench
5 flag 处理命令行参数

测试：
默认参数 500个客户端，发送1000个大小为100的包 用时8秒

异步提高的方案：
1、协程池
2、采用同时发送

