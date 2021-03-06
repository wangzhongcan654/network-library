/*==============================================================
#Author: curl.wang
#Description: tcp serivcie common type
#Version: v1.0.0
#LastChangeTime:2018-12-07
==============================================================*/
package netcommon

import (
	"net"
	"network-library/tcpserver"
)

const(
	//默认协议头，代表真实报文数据大小
	DefaultHeadSize = 4
	//默认协议体报文大小
	DefaultPackageSize = 4* 1024
)

var (
	EventIn = func(conn *tcpserver.TcpConn,msg []byte){}
	EventOut = func(conn net.Conn,msg []byte){}
	EventConnect = func(conn *tcpserver.TcpConn){}
	EventDisconnect = func(conn *tcpserver.TcpConn){}
)

