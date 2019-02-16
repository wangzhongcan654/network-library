package tcpserver

import (
	"sync"
	"net"
	"sync/atomic"
	"time"
	"network-library/common"
)

var(
	tcpConnId uint64
	noDeadline = time.Time{}
)
const(
	defaultKeepAlivePeriod = time.Second*10

)
type TcpConn struct{
	//当前连接标识id
	id uint64
	conn net.Conn

	reader *reader
	sendMutex sync.Mutex

	writer *writer
	recvMutex sync.Mutex
	//true is close
	closeFlag bool
}

//初始化新的连接
func NewConn(conn net.Conn) *TcpConn{
	if conn == nil{
		return nil
	}
	return &TcpConn{
		id:atomic.AddUint64(&tcpConnId,1),
		conn:conn,
		reader:NewReader(conn),
		writer:new(writer),
	}
}

func (c *TcpConn) Id() uint64{
	return c.id
}

func (c *TcpConn) Conn() net.Conn{
	return c.conn
}

func (c *TcpConn) IsClosed() bool{
	return c.closeFlag
}

func (c *TcpConn) Close(){
	if c.closeFlag != true{
		c.closeFlag=true
		c.conn.Close()
		netcommon.EventDisconnect(c)
	}
}

func (c *TcpConn) Receive()(msg []byte,err error){
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()

	if msg,err=c.reader.ReadPacket() ; err!=nil{
		c.Close()
		return
	}
	return
}

func (c *TcpConn) Send(msg []byte)(n int,err error){
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	if n,err=c.writer.WritePakcage(msg);err!=nil{
		c.Close()
		return
	}
	netcommon.EventOut(c.conn,msg)
	return
}

func (c *TcpConn) SetKeepAlive(period time.Duration)(err error){
	if tc,ok:=c.conn.(*net.TCPConn);ok{
		if err=tc.SetKeepAlive(true);err!=nil{
			return
		}
		if err=tc.SetKeepAlivePeriod(period);err!=nil{
			return
		}
	}

	return
}

func (c *TcpConn)SetDeadline(timeout time.Duration)(err error){
	if tc,ok:=c.conn.(*net.TCPConn);ok{
		if timeout!=0{
			err=tc.SetDeadline(time.Now().Add(timeout))
		}else{
			err=tc.SetDeadline(noDeadline)
		}
	}
	return
}

func (c *TcpConn)SetReadDeadline(timeout time.Duration)(err error){
	if tc,ok:=c.conn.(*net.TCPConn);ok{
		if timeout!=0{
			err=tc.SetReadDeadline(time.Now().Add(timeout))
		}else{
			err=tc.SetReadDeadline(noDeadline)
		}
	}
	return
}

func (c *TcpConn) SetWriteDeadline(timeout time.Duration)(err error){
	if tc,ok:=c.conn.(*net.TCPConn);ok{
		if timeout!=0{
			tc.SetWriteDeadline(time.Now().Add(timeout))
		}else{
			tc.SetWriteDeadline(noDeadline)
		}
	}
	return
}