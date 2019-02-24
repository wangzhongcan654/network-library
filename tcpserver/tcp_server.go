package tcpserver

import (
	"fmt"
	"net"
	"network-library/common"
	"sync"
	"sync/atomic"
)

//tcp服务器连接管理对象
type tcpServer struct{
	listener net.Listener

	maxConnectionId uint64
	connections map[uint64]*TcpConn
	connMutex sync.Mutex

	stopFlag int32
	stopWait sync.WaitGroup

	currConn uint64
}

/*
*@NewTcpServer: 创建一个tcpserver
*@listener：监听对象
*@return: tcpserver对象
*/
func NewTcpServer(listener net.Listener) *tcpServer{
	if listener!=nil{
		return &tcpServer{
			listener:listener,
			connections: make(map[uint64]*TcpConn),
			currConn:0,
		}
	}
	return nil
}


/*
*@beginServe: 创建一个tcpserver
*@return: 服务开启过程的错误信息
*/
func (s *tcpServer) beginServe()(err error){
	fmt.Printf("server start service, listen to %s",s.listener.Addr())
	for{
		conn,err:=s.listener.Accept()
		if err!=nil{
			fmt.Printf("Accept failed: %v",err)
			break;
		}
		connection,_:=s.newConn(conn)

		fmt.Println("New client connection [%s -> %v]",conn.RemoteAddr(),conn.LocalAddr())

		netcommon.EventConnect(connection)


	}
	return
}

/*
*@newConn: 创建一个tcpserver
*@conn：监听对象
*@return: 返回单个连接TcpConn对象，出错并返回错误信息
*/

func (s *tcpServer) newConn(conn net.Conn)(c *TcpConn,err error){
	c=NewConn(conn)
	if err=c.SetKeepAlive(defaultKeepAlivePeriod);err!=nil{
		return
	}
	s.addConn(c)
	return
}

/*
*@addConn: 添加一个连接到连接管理表
*@conn：单个连接对象
*@return: no
*/
func (s *tcpServer) addConn(conn *TcpConn){
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	s.connections[conn.id]=conn
	s.stopWait.Add(1)
}

/*
*@serveConnection: 添加一个连接到连接管理表
*@conn：单个连接对象
*@return: no
*/
func (s *tcpServer) serveConnection(conn *TcpConn){
	for{
		req,err:=conn.Receive()
		if err!=nil{
			fmt.Printf("conntion [%s -> %v] is closed",conn.conn.RemoteAddr(),conn.conn.LocalAddr())
			s.deleteConn(conn)
			return
		}

		go netcommon.EventIn(conn,req)
	}
}

/*
*@lenConnection: 获取当前tcp服务器连接数
*@return: 当前服务器连接数
*/
func (s *tcpServer) lenConnection() int{
	return len(s.connections)
}


/*
*@stop: 停止当前服务器服务
*@return: no
*/
func (s *tcpServer) stop() bool{
	if atomic.CompareAndSwapInt32(&s.stopFlag,0,1){
		s.listener.Close()
		s.closeConn()
		s.stopWait.Wait()
		return true
	}
	return false
}

/*
*@deleteConn: 关闭单个连接
*@conn：单个连接对象
*@return: no
*/
func (s *tcpServer) deleteConn(conn *TcpConn){
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	delete(s.connections,conn.id)
	s.stopWait.Done()
}

/*
*@closeConn: 关闭当前服务器所有连接
*@return: no
*/
func (s *tcpServer) closeConn(){
	for _,conn:=range s.connections{
		conn.Close()
	}
}