package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func StartHttpService(port string) error{
	localAddress, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%v",port))//定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", localAddress)       //在刚定义好的地址上进监听请求。
	if err != nil {
		return err
	}
	defer func() { //担心return之前忘记关闭连接，因此在defer中先约定好关它。
		tcpListener.Close()
	}()
	log.Printf("[Info]:start listen")
	for {
		var conn, err = tcpListener.AcceptTCP() //接受连接。
		if err != nil {
			return err
		}
		var remoteAddr = conn.RemoteAddr() //获取连接到的对像的IP地址。
		log.Printf("[Info]:connect to %v",remoteAddr.String())
		bys, err := ioutil.ReadAll(conn) //读取对方发来的内容。
		if err != nil {
			return err
		}
		log.Printf("[Info]:get message :%v",bys)

		// 服务转发
		err = forward(bys)
		if err != nil {
			log.Printf("[err]:%v",err)
		}
		conn.Write([]byte("hello, Nice to meet you, my name is SongXingzhu")) //尝试发送消息。
		conn.Close()
	}

}

func forward(data []byte) error {
	log.Printf("[Info]:start forward")
	rev := make(map[string]interface{})
	var err error

	err = json.Unmarshal(data,&rev)
	if err != nil {
		return err
	}
	log.Printf("[Info]:unmarshal ok")
	log.Printf("[Info]:message name is %v", rev["MessageName"])
	switch rev["MessageName"].(string) {
	case "RecvReq":
		err = RecvVerifyReq(data)
	}

	if err != nil {
		return err
	}
	return nil
}

