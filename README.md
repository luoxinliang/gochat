本示例小程序主要练习go的rpc，go程。

聊天小程序分为两部分：server，client。分开运行：

  1.运行服务器
    cd gochat_server/src/ 
    go run main.go    
    按提示，输入：
    command> start
    则启动服务器了
  
  2.运行客户端（新的终端下）
    cd gochat_client/src/ 
    go run main.go
  
    按提示，输入：
    command> start 7777  
    则启动了一个客户端
    command> login tom
    登录，用户名为 tom
  
    打开新的shell，启动新的客户端，登录更多用户
    如 command>start 9999
       command>login jack
  
  于是，就可以聊天了。比如，tom向jack说句“hello”，在jack的客户端会收到消息。
  
    command>sendTo jack hello
  


  
  