package server

import (
	"fmt"
	"net/netip"
	"sync"
	"time"
	"vchat/user"

	netreactors "github.com/lance-e/net-reactors"
)

type Server struct {
	Server_ *netreactors.TcpServer
	//在线用户map表;key为name,value为User对象
	OnlineMap map[string]*user.User
	//读写锁
	MapLock sync.RWMutex
	//用户与tcp连接的map
	ConnToUser map[*netreactors.TcpConnection]*user.User
}

func NewServer(ev *netreactors.EventLoop, ip string, port int, name string, goroutineNum int) *Server {
	addport, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic("NewServer:ip or port wrong\n")
	}
	S := Server{
		Server_:    netreactors.NewTcpServer(ev, &addport, name),
		OnlineMap:  make(map[string]*user.User),
		MapLock:    sync.RWMutex{},
		ConnToUser: make(map[*netreactors.TcpConnection]*user.User),
	}
	S.Server_.SetConnectionCallback(S.NewConnection)
	S.Server_.SetMessageCallback(S.Broadcast)
	S.Server_.SetGoroutineNum(goroutineNum)

	return &S
}

// 建立连接的回调
func (S *Server) NewConnection(conn *netreactors.TcpConnection) {
	if conn.Connected() {
		fmt.Printf("Server:connection established\n")
		//用户上线
		S.Online(conn)
	} else {
		fmt.Printf("Server:connection destroyed\n")
		//用户下线
		S.Offline(conn)
	}

}

// 接收消息的回调:广播消息(由哪个用户发起和发送的信息)
func (S *Server) Broadcast(conn *netreactors.TcpConnection, buffer *netreactors.Buffer, t time.Time) {
	msg := fmt.Sprintf("[%s][%s] %s:%s", conn.PeerAddr().String(), t.Format("2006-01-02 15:04:05"), conn.Name(), buffer.RetrieveAllString())
	S.MapLock.Lock()
	for _, cli := range S.OnlineMap {
		cli.C <- msg
	}
	S.MapLock.Unlock()
}

func (S *Server) Start() {
	S.Server_.Start()
}

func (S *Server) Online(conn *netreactors.TcpConnection) {
	u := user.NewUser(conn)

	msg := fmt.Sprintf("[%s][%s] %s:%s", u.Conn.PeerAddr().String(), time.Now().Format("2006-01-02 15:04:05"), u.Name, "已上线\n")

	S.MapLock.Lock()
	//广播用户上线
	for _, cli := range S.OnlineMap {
		cli.C <- msg
	}
	//将上线用户加入online表和ConnToUser表中
	S.OnlineMap[u.Name] = u
	S.ConnToUser[conn] = u

	S.MapLock.Unlock()

	u.Conn.Send([]byte("welcome server\n"))
}

func (S *Server) Offline(conn *netreactors.TcpConnection) {
	u := S.ConnToUser[conn]
	msg := fmt.Sprintf("[%s][%s] %s:%s", u.Conn.PeerAddr().String(), time.Now().Format("2006-01-02 15:04:05"), u.Name, "已离线\n")

	S.MapLock.Lock()
	//将离线用户从online表中删除
	delete(S.OnlineMap, u.Name)
	delete(S.ConnToUser, conn)

	//广播用户离线
	for _, cli := range S.OnlineMap {
		cli.C <- msg
	}

	S.MapLock.Unlock()
}
