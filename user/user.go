package user

import (
	"fmt"

	netreactors "github.com/lance-e/net-reactors"
)

type User struct {
	Name string
	Addr string
	C    chan string                //跟用户绑定
	Conn *netreactors.TcpConnection //用户用于通信
}

func NewUser(connection *netreactors.TcpConnection) (user *User) {
	user = &User{
		Name: connection.Name(),
		Addr: connection.LocalAddr().String(),
		C:    make(chan string),
		Conn: connection,
	}
	go user.TransferMessage()
	return
}

func (u *User) TransferMessage() {
	for {
		msg := <-u.C
		fmt.Printf("user[%s] get message:%s\n", u.Name, msg)
		u.Conn.Send([]byte(msg))
	}
}
