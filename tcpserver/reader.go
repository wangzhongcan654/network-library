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

type reader struct{
	buf []byte
	read io.Reader
}

/*
*@NewReader:返回要给reader对象的指针
*@r：要reader的io对象
*/
func NewReader(r io.Reader)*reader{
	return &reader{
		buf:make([]byte,netcommon.DefaultPackageSize),
		read:r,
	}
}

/*
*@ReadPacket:读取io中的一个完整的数据包
*@r：要reader的io对象
*@return：返回一个silce，这个silce存储读取到的数据
*/
func (r *reader)ReadPacket()(packet []byte,err error){
	sizePacket,err:=r.readHead()
	if err!=nil{
		return
	}
	_,err=r.readBody(sizePacket)
	if err!=nil{
		return
	}
	packet=make([]byte,0,sizePacket)
	packet=append(packet,r.buf[:sizePacket]...)
	return
}

/*
*@readHead:读取协议头
*@r：要reader的io对象
*@return： 返回当前读取的packet的body大小
*/
func (r *reader)readHead()(len int,err error){
	n,err:=r.readUint32BigEndian()
	if err!=nil{
		return 0,err
	}
	return int(n),nil
}

/*
*@readUint32BigEndian:读取头部信息，并将网络字节序转换为本地字节序
*@r：要reader的io对象
*@return：返回当前转换后的头部信息（完整数据包大小）
*/
func (r *reader)readUint32BigEndian()(n uint32,err error){
	_,err=io.ReadFull(r.read,r.buf[:netcommon.DefaultHeadSize])
	if err!=nil{
		return 0,err
	}
	n=binary.BigEndian.Uint32(r.buf[:netcommon.DefaultHeadSize])
	return
}

/*
*@readBody:读取协议头
*@r：要reader的io对象
*@len：要从io读取的字节大小
*@return：返回读取到的字节大小
*/
func (r*reader)readBody(len int)(n int,err error){
	//这种增加缓冲区的方式有没有可能造成内存泄漏
	if len>netcommon.DefaultPackageSize{
		r.buf=make([]byte,len)
	}
	return io.ReadFull(r.read,r.buf[:len])
}