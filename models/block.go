package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"
)

var HTTPPort string
var GRPCPort string
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
		return errors.New("Block is not to be recovered")
	}
	b.Token = "safe token"

	recoverNum := (len(BlockChain) - 1) / 3

	recoverList, err := GetTopNDist(b.Index, recoverNum)
	if err != nil {
		return err
	}
	log.Printf("[Info]:Get top dist successful\n")
	log.Printf("[Data]:RecoverList %v\n", recoverList)
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
		log.Printf("[Info]:Send msg to %v:%v", BlockChain[rq.To].Addr.Ip, BlockChain[rq.To].Addr.Port)
	}
	return nil
}

// 0+1 call
func (b *Block) SendRecvVerifyReq(req RecoverReq) error {
	if b.Type == 2 {
		return errors.New("Block is not allowed to send recover verify request")
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
		log.Printf("[Info]:Send msg to %v:%v", BlockChain[rq.To].Addr.Ip, BlockChain[rq.To].Addr.Port)
		if err != nil {
			return err
		}
	}
	return nil
}

// 1 call
func (b *Block) SendRecvRes(req []RecoverVerifyReq) error {
	if b.Type != 1 {
		return errors.New("Block is not recover node")
	}
	if len(req)*3 < len(req[0].BlockList)*2 {
		return errors.New("Do not get enough verify replay")
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
	log.Printf("[Info]:Send msg to %v:%v", BlockChain[rq.To].Addr.Ip, BlockChain[rq.To].Addr.Port)
	if err != nil {
		return err
	}
	return nil
}

func (b *Block) SendRecvMd5Req(res []RecoverRes) error {
	if len(res) != len(res[0].BlockList) {
		return errors.New("Do not get enough recover replay")
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
		var m1 int
		var m2 int
		var m3 int
		sliceTmp := make(map[int]int)
		for {
			for key := range map1 {
				if _, ok := sliceTmp[key]; ok {
					continue
				} else {
					sliceTmp[key] = key
					m1 = key
					//delete(map1, key)
					break
				}
			}
			for key := range map2 {
				if _, ok := sliceTmp[key]; ok {
					continue
				} else {
					sliceTmp[key] = key
					m2 = key
					//delete(map2, key)
					break
				}
			}
			for key := range map3 {
				if _, ok := sliceTmp[key]; ok {
					continue
				} else {
					sliceTmp[key] = key
					m3 = key
					//delete(map3, key)
					break
				}
			}
			if len(sliceTmp) == 3 {
				delete(map1, m1)
				delete(map2, m2)
				delete(map3, m3)
				break
			} else {
				sliceTmp = make(map[int]int)
			}
		}

		sliceSize := res[0].Size / int64(sliceNum)

		dataSlice := make([]DataSlice, 0)
		for key := range sliceTmp {
			if key != len(res)-1 {
				dataSlice = append(dataSlice,DataSlice{
					Offset: int64(key) * sliceSize,
					Size:   sliceSize,
				})
			} else {
				dataSlice = append(dataSlice,DataSlice{
					Offset: int64(key) * sliceSize,
					Size:   res[0].Size - int64(key)*sliceSize,
				})
			}
		}

		rq := &RecoverMd5Req{
			BaseMessage: BaseMessage{
				MessageId:   GetRandomId(),
				From:        b.Index,
				To:          res[0].BlockList[i],
				MessageName: "SendRecvMd5Req",
				Timestamp:   time.Now().String(),
			},
			Slices: dataSlice,
		}
		err := SendMsg(BlockChain[rq.To].Addr, rq)
		log.Printf("[Info]:Send msg to %v:%v", BlockChain[rq.To].Addr.Ip, BlockChain[rq.To].Addr.Port)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) SendRecvMd5Res(req RecoverMd5Req) error {
	data := []byte(b.Data)
	for idx, e := range req.Slices {
		dataSlc := data[e.Offset : e.Offset+e.Size]
		req.Slices[idx].Md5 = fmt.Sprintf("%v", md5.Sum(dataSlc))
	}

	rq := &RecoverMd5Res{
		BaseMessage: BaseMessage{
			MessageId:   GetRandomId(),
			From:        b.Index,
			To:          req.From,
			MessageName: "SendRecvMd5Res",
			Timestamp:   time.Now().String(),
		},
		Slices: req.Slices,
	}
	err := SendMsg(BlockChain[rq.To].Addr, rq)
	log.Printf("[Info]:Send msg to %v:%v", BlockChain[rq.To].Addr.Ip, BlockChain[rq.To].Addr.Port)
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
