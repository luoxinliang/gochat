/**
 * Author: luoxinliang
 * Email:luoxinliang.zh@gmail.com
 * 用Go，就够了
 * Date: 12-9-19
 * Time: 上午10:49
 */
package server

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
	SERVER_PORT = "9180"            //服务器默认端口 9180
	SERVER_ADDR = SERVER_IP + ":" + SERVER_PORT
)

type Message struct {
	From    *User
	To      *User
	Content string
}

type Reply struct {
	StateCode   int32
	Content     string
	Error       string
}

func sendMsgToUser(message *Message) (reply *Reply)  {
	toUserName := message.To.UserName
	toUser := getUser(toUserName)
	if toUser == nil {
		reply.StateCode = 503
		reply.Error = "sendMsgToUser fail: toUser not exsist!"
		return reply
	}
	method := "User.ShowMessage"
	replyCall := callClientUser(toUser, method, message, &Reply{})
	result := <-replyCall.Done
	return result.Reply.(*Reply)
}

//服务器向在线用户发信息
func sendMsg(args []string) (reply *Reply)  {
	if len(args) != 3 {
		fmt.Println("Need 3 args: <method> <to-user-name> <message-content>")
		reply.StateCode = 503
		reply.Error = "Args error"
		return reply
	}
	toUserName := args[1]
	messageContent := args[2]
	toUser := getUser(toUserName)
	if toUser == nil {
		reply.StateCode = 503
		reply.Error = "SendMsg failed: toUser not exists..."
		return reply
	}
	message := &Message{To:toUser, Content:messageContent}
	return sendMsgToUser(message)
}

func startAccept(args []string) (reply *Reply) {
	if len(args) != 1 {
		fmt.Println("Need one args <method>")
		reply.StateCode = 503
		reply.Error = "Args error"
		return reply
	}
	l, e := net.Listen(TCP, SERVER_ADDR)
	if e != nil {
		fmt.Println("Listen error", e)
		reply.StateCode = 503
		reply.Error = "Listen error"
		return reply
	}

	fmt.Println("Server: listened on PORT  ", SERVER_PORT)

	go func() {
		for {
			conn, e := l.Accept()
			if e != nil {
				fmt.Println("Accept err", e)
				conn.Close()
				continue
			}
			fmt.Println("Accept From:", conn.RemoteAddr())
			if conn != nil {
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

//启动服务器，默认端口9180。必须在客户端启动前启动服务器
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

func callClientUser(u *User, method string, args interface{}, reply interface{}) (*rpc.Call) {
	return callClient(u.CurIp, u.CurPort, method, args, reply)
}

func callClient(ip string, port string, method string, args interface{}, reply interface{}) (*rpc.Call) {
	dialAddr := ip + ":" + port
	fmt.Println("dialAddr:", dialAddr)
	client, err := jsonrpc.Dial(TCP, dialAddr)
	if err != nil {
		fmt.Println("Dial error...")
		return nil
	}
	return client.Go(method, args, reply, nil)
}

func Help(args []string) *Reply {
	fmt.Println(`
Commands:
	start
	sendMsg <to-user-name> <content>
	help<h>
		`)
	return nil
}

func getCommandHandler() map[string]func (args []string) *Reply {
	return map[string] func ([]string) *Reply {
		"start":startAccept,
		"sendMsg":sendMsg,
		"help<h>":Help,
		"h":Help,
	}
}

