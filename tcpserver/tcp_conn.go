/*==============================================================
#Author: curl.wang
#Description: tcpconn  single conn
#Version: v1.0.0
#LastChangeTime:2019-02-24
==============================================================*/
package tcpserver

import (
	"network-library/common"
	"sync"
	"net"
	"sync/atomic"
	"time"

)

var(
	tcpConnId uint64
	noDeadline = time.Time{}
)
const(
	defaultKeepAlivePeriod = time.Second*10

)

//单个tcp连接管理对象
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

/*
*@NewConn: 构造一个tcp连接对象
*@conn：连接对象
*@return: 单个连接管理对象
*/
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

/*
*@Id:返回当前连接id
*@return:当前连接id
*/
func (c *TcpConn) Id() uint64{
	return c.id
}

/*
*@Conn: 获取当前连接对象的网络连接
*@return:网络连接对象
*/
func (c *TcpConn) Conn() net.Conn{
	return c.conn
}

/*
*@IsClosed:判断当前连接是否已经关闭
*@return: false关闭
*/
func (c *TcpConn) IsClosed() bool{
	return c.closeFlag
}

/*
*@Close:关闭当前网络连接
*@return: no
*/
func (c *TcpConn) Close(){
	if c.closeFlag != true{
		c.closeFlag=true
		c.conn.Close()
		netcommon.EventDisconnect(c)
	}
}

/*
*@Receive: 接收包
*@return: msg接收到的包数据，接收出错返回具体错误信息
*/
func (c *TcpConn) Receive()(msg []byte,err error){
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()

	if msg,err=c.reader.ReadPacket() ; err!=nil{
		c.Close()
		return
	}
	return
}

/*
*@Send:发送数据包
*@msg：要发送的数据
*@return: n：发送成功的字节数，err：失败返回具体的错误信息
*/
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

/*
*@SetKeepAlive:设置当前连接的保持时间
*@period： 保活时间
*@return:设置错误返回具体出错信息
*/
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

/*
*@SetDeadline:构造一个writer对象，并返回对象的指针
*@timeout：要writer的io对象
*@return:对象指针
*/
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

/*
*@SetReadDeadline:构造一个writer对象，并返回对象的指针
*@timeout：要writer的io对象
*@return:对象指针
*/
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

/*
*@SetWriteDeadline:构造一个writer对象，并返回对象的指针
*@timeout：要writer的io对象
*@return:对象指针
*/
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