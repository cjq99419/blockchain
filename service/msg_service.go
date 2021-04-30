package service

import (
	"blockchain/models"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

var RecvVerifyList []models.RecoverVerifyReq
var RecvResList []models.RecoverRes

func ServeRecvReq(data []byte) error {
	req := &models.RecoverReq{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	idx,err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}

	log.Printf("[Info]:unmarshal ok %v", req)
	err = models.BlockChain[idx].SendRecvVerifyReq(*req)
	if err != nil {
		return err
	}
	return nil
}

func ServeRecvVerifyReq(data []byte) error {
	if RecvVerifyList == nil {
		RecvVerifyList = make([]models.RecoverVerifyReq, 0)
	}

	req := &models.RecoverVerifyReq{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}
	RecvVerifyList = append(RecvVerifyList, *req)
	log.Printf("[Info]:get replay : %v", len(RecvVerifyList) - 1)
	return nil
}

func TryServeRecvVerifyReq() error {
	time.Sleep(1 * time.Second)
	return recvRes()
}

func recvRes() error{
	idx,err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}
	err = models.BlockChain[idx].SendRecvRes(RecvVerifyList)
	if err != nil {
		return err
	}
	log.Printf("[Info]:recv res")
	RecvVerifyList = make([]models.RecoverVerifyReq, 0)
	return nil
}

func ServeRecvRes(data []byte) error {
	if RecvResList == nil {
		RecvResList = make([]models.RecoverRes, 0)
	}

	res := &models.RecoverRes{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return err
	}

	idx,err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}

	if len(res.BlockList) == len(RecvResList) {
		err = models.BlockChain[idx].SendRecvMd5Req(RecvResList)
		if err != nil {
			return nil
		}
		RecvResList = make([]models.RecoverRes, 0)
	}
	return nil
}