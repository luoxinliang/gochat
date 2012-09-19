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
	"net/rpc"
)

type User struct {
	Id        int32
	UserName  string
	CurIp     string
	CurPort   string
}

func (user *User)ShowMessage(message *Message,reply *Reply) error {
	fmt.Printf("%s recieved a message from %s.Content:%s",message.To.UserName,message.From.UserName,message.Content)
	fmt.Println()
	reply.StateCode = 200
	reply.Content = "Reciever recieved the message."
	return nil
}

func RegisterRPC() {
	rpc.Register(new(User))
}



