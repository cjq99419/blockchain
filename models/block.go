package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"
)

var Port string
var Index string

var BlockChain []Block

type Block struct {
	Id string

	Index int
	//create time
	Timestamp string
	Data      string
	//current hash
	Hash    string
	PreHash string

	//经度
	Lon string
	//纬度
	Lat string

	//network addr
	Addr Address

	//token
	Token string

	//0:other block, 1:recover block, 2:to be recovered block
	Type int
}

// 2 call
func (b *Block) SendRecvReq() error {
	if b.Type != 2 {
		return errors.New("block is not to be recovered")
	}
	b.Token = "safe token"

	recoverNum := (len(BlockChain) - 1) / 3

	recoverList, err := GetTopNDist(b.Index, recoverNum)
	if err != nil {
		return err
	}
	log.Printf("[Info]:get top dist successful\n")
	log.Printf("[Data]:recoverList %v\n",recoverList)
	for i := 0; i < len(BlockChain); i++ {
		if i == b.Index {
			continue
		}
		rq := &RecoverReq{
			BaseMessage: BaseMessage{
				MessageId:   GetRandomId(),
				From:        b.Index,
				To:          BlockChain[i].Index,
				MessageName: "SendRecvReq",
				Timestamp:   time.Now().String(),
			},
			Token:     b.Token,
			BlockList: recoverList,
		}

		err = SendMsg(BlockChain[rq.To].Addr, rq)
		if err != nil {
			return err
		}
		log.Printf("[Info]:send msg to %v:%v",BlockChain[rq.To].Addr.Ip,BlockChain[rq.To].Addr.Port)
	}
	return nil
}

// 0+1 call
func (b *Block) SendRecvVerifyReq(req RecoverReq) error {
	if b.Type == 2 {
		return errors.New("block is not allowed to send recover verify request")
	}
	flag := true
	if req.Token != "safe token" {
		flag = false
	}



	for _, idx := range req.BlockList {
		rq := &RecoverVerifyReq{
			BaseMessage: BaseMessage{
				MessageId:   GetRandomId(),
				From:        b.Index,
				To:          BlockChain[idx].Index,
				MessageName: "SendRecvVerifyReq",
				Timestamp:   time.Now().String(),
			},
			ReqId:     req.MessageId,
			isVerify:  flag,
			RecvBlock: req.From,
			BlockList: req.BlockList,
		}

		BlockChain[idx].Type = 1

		err := SendMsg(BlockChain[rq.To].Addr, rq)
		if err != nil {
			return err
		}
	}
	return nil
}

// 1 call
func (b *Block) SendRecvRes(req []RecoverVerifyReq) error {
	if b.Type != 1 {
		return errors.New("block is not recover node")
	}
	if len(req)*3 < len(req[0].BlockList)*2 {
		return errors.New("do not get enough verify replay")
	}

	rq := &RecoverRes{
		BaseMessage: BaseMessage{
			MessageId:   GetRandomId(),
			From:        b.Index,
			To:          req[0].RecvBlock,
			MessageName: "SendRecvRes",
			Timestamp:   time.Now().String(),
		},
		ReqId:     req[0].ReqId,
		Size:      int64(len(b.Data)),
		BlockList: req[0].BlockList,
	}

	err := SendMsg(BlockChain[rq.To].Addr, rq)
	if err != nil {
		return err
	}
	return nil
}

func (b *Block) SendRecvMd5Req(res []RecoverRes) error {
	if len(res) != len(res[0].BlockList) {
		return errors.New("do not get enough recover replay")
	}
	sliceNum := len(res[0].BlockList)
	map1 := make(map[int]int)
	map2 := make(map[int]int)
	map3 := make(map[int]int)
	for i := 0; i < sliceNum; i++ {
		map1[i] = i
		map2[i] = i
		map3[i] = i
	}

	for i := 0; i < len(res); i++ {
		sliceTmp := make(map[int]int)
		for key := range map1 {
			if _, ok := sliceTmp[key]; ok {
				continue
			} else {
				sliceTmp[key] = key
				delete(map1, key)
				break
			}
		}
		for key := range map2 {
			if _, ok := sliceTmp[key]; ok {
				continue
			} else {
				sliceTmp[key] = key
				delete(map2, key)
				break
			}
		}
		for key := range map3 {
			if _, ok := sliceTmp[key]; ok {
				continue
			} else {
				sliceTmp[key] = key
				delete(map3, key)
				break
			}
		}

		sliceSize := res[0].Size / int64(sliceNum)

		var dataSlice [3]DataSlice
		var j int
		for key := range sliceTmp {
			if key != len(res)-1 {
				dataSlice[j] = DataSlice{
					offset: int64(key) * sliceSize,
					size:   sliceSize,
				}
			} else {
				dataSlice[j] = DataSlice{
					offset: int64(key) * sliceSize,
					size:   res[0].Size - int64(key)*sliceSize,
				}
			}
			j++
		}

		rq := &RecoverMd5Req{
			BaseMessage: BaseMessage{
				MessageId:   GetRandomId(),
				From:        b.Index,
				To:          res[0].BlockList[i],
				MessageName: "SendRecvMd5Req",
				Timestamp:   time.Now().String(),
			},
			dataSlice: dataSlice,
		}
		err := SendMsg(BlockChain[rq.To].Addr, rq)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) SendRecvMd5Res(req RecoverMd5Req) error {
	data := []byte(b.Data)
	for idx, e := range req.dataSlice {
		dataSlc := data[e.offset : e.offset+e.size]
		req.dataSlice[idx].md5 = fmt.Sprintf("%v", md5.Sum(dataSlc))
	}

	rq := &RecoverMd5Res{
		BaseMessage: BaseMessage{
			MessageId:   GetRandomId(),
			From:        b.Index,
			To:          req.From,
			MessageName: "SendRecvMd5Res",
			Timestamp:   time.Now().String(),
		},
		dataSlice: req.dataSlice,
	}
	err := SendMsg(BlockChain[rq.To].Addr, rq)
	if err != nil {
		return err
	}
	return nil
}

func (b *Block) VerifyMd5Req(to Address, req VerifyMd5Req) error {
	return nil
}

func (b *Block) VerifyMd5Res(to Address, req VerifyMd5Res) error {
	return nil
}
