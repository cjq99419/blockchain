package models

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

type Address struct {
	Ip   string
	Port string
}

type DataSlice struct {
	Offset int64
	Size   int64
	Md5    string
}

type BaseMessage struct {
	MessageId string
	From      int
	To        int

	MessageName string

	Timestamp string
}

type RecoverReq struct {
	BaseMessage
	Token     string
	BlockList []int
}

type RecoverVerifyReq struct {
	BaseMessage
	ReqId     string
	isVerify  bool
	RecvBlock int
	BlockList []int
}

type RecoverRes struct {
	BaseMessage
	ReqId     string
	Size      int64
	BlockList []int
}

type RecoverMd5Req struct {
	BaseMessage
	Slice []DataSlice
}

type RecoverMd5Res struct {
	BaseMessage
	Slice []DataSlice
}

type VerifyMd5Req struct {
	BaseMessage
}

type VerifyMd5Res struct {
	BaseMessage
	Md5 string
}

func SendMsg(addr Address, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr.Ip, addr.Port))
	if err != nil {
		return err
	}
	log.Println(string(data))
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func GetRandomId() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
