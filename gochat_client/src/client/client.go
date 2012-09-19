/**
 * Author: luoxinliang
 * Email:luoxinliang.zh@gmail.com
 * 用Go,就够了
 * Date: 12-9-19
 * Time: 下午7:46
 */
package client

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"bufio"

	"net/rpc/jsonrpc"
)

const (
	TCP         = "tcp"
	SERVER_IP   = "127.0.0.1"
	SERVER_PORT = "9180"          //服务器默认端口 9180
	SERVER_ADDR = SERVER_IP + ":" + SERVER_PORT
	CLIENT_IP   = "127.0.0.1"
)

var port string
var curUser *User

type Message struct {
	From    *User
	To      *User
	Content string
}

type Reply struct {
	StateCode   int32 //状态码 正常：200  不正常其他
	Content     string
	Error       string
}

//登录
func login(args []string) (reply *Reply) {
	reply = &Reply{}
	if len(args) != 2 {
		fmt.Println("Need 2 args <method> <username>")
		reply.StateCode = 503
		reply.Error = "Args error"
		return
	}
	if port == "" {
		fmt.Println("Client not started.Should start client first!")
		reply.StateCode = 503
		reply.Error = "Client not started  error"
		return
	}
	userName := args[1]
	user := &User{UserName:userName,CurIp:CLIENT_IP,CurPort:port}
	method := "User.Login"
	replyCall := callServer(method, user, &Reply{})
	result := <-replyCall.Done
	reply = result.Reply.(*Reply)
	if reply.StateCode == 200 {
		fmt.Println("Login success!")
		curUser = user
	} else {
		fmt.Println("Login fail! Error msg:",reply.Error)
	}
	return
}

//向用户发送消息
func sendTo(args []string) (reply *Reply) {
	reply = &Reply{}
	if len(args) != 3 {
		fmt.Println("Need 3 args <method> <to-username> <content>")
		reply.StateCode = 503
		reply.Error = "Args error"
		return
	}
	if curUser == nil {
		fmt.Println("Not logined,should login first!")
		reply.StateCode = 503
		reply.Error = "Not login error"
		return
	}
	method := "User.SendTo"
	toUserName := args[1]
	messageContent := args[2]
	toUser := &User{UserName:toUserName}
	message := &Message{To:toUser,From:curUser,Content:messageContent}
	replyCall := callServer(method, message, &Reply{})
	result := <-replyCall.Done
	reply = result.Reply.(*Reply)
	return
}

func getClientAddr() string {
	if port == "" {
		return ""
	}
	return CLIENT_IP + ":" + port
}

//启动客户端聊天，设置收信息端口
func startAccept(args []string) *Reply {
	if len(args) != 2 {
		fmt.Println("Need 2 args <method> <port>")
		return nil
	}
	port = args[1]
	clientAddr := getClientAddr()
	l, e := net.Listen(TCP, clientAddr)
	if e != nil {
		fmt.Println("Listen error", e)
		return nil
	}

	fmt.Println("Clinet: listening on PORT ", port)

	go func() {
		for {
			conn, e := l.Accept()
			if e != nil {
				fmt.Println("Accept err", e)
				conn.Close()
				continue
			}
			if conn != nil {
				fmt.Println("Accept From:", conn.RemoteAddr())
				go func() {
					rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
					conn.Close()
				}()
			}
		}
	}()

	RegisterRPC()

	return nil

}

func Start() {
	r := bufio.NewReader(os.Stdin)
	handlers := getCommandHandler()
	for {
		fmt.Print("Command>")
		b, _, _ := r.ReadLine()
		line := string(b)
		tokens := strings.Split(line, " ")
		if handler, ok := handlers[tokens[0]]; ok {
			handler(tokens)
		}
	}
}

func callServer(method string, args interface{}, reply interface{}) (*rpc.Call) {
	client, err := jsonrpc.Dial(TCP, SERVER_ADDR)
	if err != nil {
		fmt.Println("Dial error...")
		return nil
	}
	return client.Go(method, args, reply, nil)
}

func Help(args []string) *Reply {
	fmt.Println(`
Commands:
	start <port>
	login <userName>
	sendTo <to-user-name> <content>
	help<h>
	`)
	return nil
}

func getCommandHandler() map[string]func (args []string) *Reply {
	return map[string] func ([]string) *Reply {
		"start":startAccept,
		"login":login,
		"sendTo":sendTo,
		"help<h>":Help,
		"h":Help,
	}
}



