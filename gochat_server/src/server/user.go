/**
 * Author: luoxinliang
 * Email:luoxinliang.zh@gmail.com
 * 用Go，就够了
 * Date: 12-9-19
 * Time: 上午10:49
 */

package server

import (
	"errors"
	"fmt"
	"net/rpc"
)

type User struct {
	Id        int32
	UserName  string
	CurIp     string
	CurPort   string
}

var OnlineUsers = make(map[string] *User)

func userLogin(u *User) {
	OnlineUsers[u.UserName] = u
	fmt.Sprintf("%s logined,Ip: %s,Port: %s",u.UserName,u.CurIp,u.CurPort)
	fmt.Sprintf("Now total %d users online",len(OnlineUsers))
}

func getUser(userName string) *User {
	u, ok := OnlineUsers[userName]
	if ok {
		return u
	}
	return nil
}

func (user *User) Login(u *User, reply *Reply) error {
	uu := getUser(u.UserName)
	if uu != nil {
		reply.StateCode = 503
		reply.Error = "Username was used,change another one!"
		return errors.New("Username was used.Login failed")
	}
	userLogin(u)
	reply.StateCode = 200
	reply.Content = "Login Success!"
	return nil
}

func (user *User) Logout(userName string, reply *Reply) error {
	delete(OnlineUsers, userName)
	if getUser(userName) != nil {
		reply.StateCode = 503
		reply.Error = "Logout fail!"
		return errors.New("Logout fail!")
	}
	reply.StateCode = 200
	reply.Content = "Logout success !"
	return nil
}

func (user *User)SendTo(message *Message,reply *Reply) error {
	fmt.Println(message.From,message.To,message.Content)
	reply = sendMsgToUser(message)
	if reply.StateCode != 200 {
		return errors.New(reply.Error)
	}
	return nil
}

func RegisterRPC() {
	rpc.Register(new(User))
}

