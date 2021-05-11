package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

var RecvNum int

func StartHttpService(port string) error {
	localAddress, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%v", port)) //定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", localAddress)                        //在刚定义好的地址上进监听请求。
	if err != nil {
		return err
	}
	defer func() { //担心return之前忘记关闭连接，因此在defer中先约定好关它。
		tcpListener.Close()
	}()
	log.Printf("[Info]:Start tcp listen to 127.0.0.1:%v", port)
	for {
		var conn, err = tcpListener.AcceptTCP() //接受连接。
		if err != nil {
			return err
		}
		var remoteAddr = conn.RemoteAddr() //获取连接到的对像的IP地址。
		log.Printf("[Info]:Connect to %v", remoteAddr.String())
		bys, err := ioutil.ReadAll(conn) //读取对方发来的内容。
		if err != nil {
			return err
		}
		log.Printf("[Info]:Get message from %v", remoteAddr.String())
		log.Printf("[Info]:Message is %v", string(bys))
		// 服务转发
		err = forward(bys)
		if err != nil {
			log.Printf("[Error]:%v", err)
		}
		conn.Close()
	}

}

func forward(data []byte) error {
	log.Printf("[Info]:Start service forward")
	rev := make(map[string]interface{})
	var err error

	err = json.Unmarshal(data, &rev)
	if err != nil {
		return err
	}
	log.Printf("[Info]:Message name is %v", rev["MessageName"])
	switch rev["MessageName"].(string) {
	case "SendRecvReq":
		err = ServeRecvReq(data)
	case "SendRecvVerifyReq":
		if len(RecvVerifyList) == 0 {
			go TryServeRecvVerifyReq()
		}
		err = ServeRecvVerifyReq(data)
	case "SendRecvRes":
		err = ServeRecvRes(data)
	case "SendRecvMd5Req":
		err = ServeRecvMd5Req(data)
	case "SendRecvMd5Res":
		err = ServeRecvMd5Res(data)
	default:
		log.Printf("[Warn]:Unknown service")
	}
	log.Printf("[Info]:Service forward finished")
	if err != nil {
		return err
	}
	return nil
}
