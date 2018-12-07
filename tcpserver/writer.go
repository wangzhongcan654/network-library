/*==============================================================
#Author: curl.wang
#Description: tcp serivcie  writer
#Version: v1.0.0
#LastChangeTime:2018-12-07
==============================================================*/
package tcpserver

import (
	"encoding/binary"
	"io"
	"network-library/common"
)

type writer struct{
	write io.Writer
	buf []byte
}

/*
*@newWriter:构造一个writer对象，并返回对象的指针
*@w：要writer的io对象
*@return:对象指针
*/
func NewWriter(w io.Writer) *writer{
	return &writer{
		write:w,
		buf:make([]byte,netcommon.DefaultHeadSize),
	}
}

/*
*@writePakcage：发送数据包
*@package: 要写入的数据包
*@return:返回真实写入的字节大小
*/
func (w *writer) WritePakcage(packet []byte)(n int,err error){
	//发送头部
	n,err=w.writeHead(len(packet))
	if err!=nil{
		return 0,err
	}
	//发送包体
	return w.writerBody(packet)
}

/*
*@writeHead:发送数据包头部
*@args:len:发送数据包body的长度
*@return:返回成功写入的字节数（头部信息所占字节数）
*/
func (w *writer)writeHead(len int)(int,error){
	return w.writerUint32BigEndian(uint32(len))
}

/*
*@writerUint32BigEndian：将要发送的信息转换为网络字节序，并发送
*@v：要发送package的大小
@return:返回真实发送数据的字节大小
*/
func (w *writer)writerUint32BigEndian(v uint32)(n int,err error){
	binary.BigEndian.PutUint32(w.buf[:netcommon.DefaultHeadSize],v)
	return w.write.Write(w.buf[:netcommon.DefaultHeadSize])
}


/*
*@writerBody:发送package的body
*@body:要发送的body
*@return：返回真实写入的字节数
*/
func (w *writer)writerBody(body []byte)(n int,err error){
	return w.write.Write(body)
}
