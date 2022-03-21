package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

//封包、拆包的具体模块
type DataPack struct{}

// 封包、拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度的方法
func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32 4字节
	//ID uint32 4字节
	return 8
}

// Pack 封包方法
// dataLen+msgId+data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//二进制序列化
	//将DataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法
// 只需要将包的head信息读出，之后再根据head信息里的data长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制流读取数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	//只解压head信息，得到dataLen和msgId
	msg := &Message{}
	//读dataLen
	//TODO 这里是默认读4个字节吗
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen是否超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("message length exceed")
	}
	return msg, nil
}
