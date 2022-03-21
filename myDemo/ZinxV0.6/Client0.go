package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
  模拟客户端
*/
func main() {
	fmt.Println("client start")

	time.Sleep(1 * time.Second)

	//1 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client0 start error! exit...")
		return
	}
	for {
		//发送封包msg消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error: ", err)
			return
		}
		//服务器回复一个msg数据，拆包
		//1 先读取流中的head部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Read head error: ", err)
			break
		}
		//2 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Client unpack msgHead error: ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//3 根据dataLen再次读取
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("Client read msg error: ", err)
				break
			}
			fmt.Println("------> Receive Server Msg: ID = ", msg.GetMsgId(), " data = ", string(msg.GetData()))
		}
		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
