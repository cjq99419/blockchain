package service

import (
	"blockchain/models"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

var RecvVerifyList []models.RecoverVerifyReq
var RecvResList []models.RecoverRes
var RecvMd5ResList []models.RecoverMd5Res

func ServeRecvReq(data []byte) error {
	req := &models.RecoverReq{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	idx, err := strconv.Atoi(models.Index)
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
	log.Printf("[Info]:get replay : %v", len(RecvVerifyList)-1)
	return nil
}

func TryServeRecvVerifyReq() error {
	time.Sleep(1 * time.Second)
	return recvRes()
}

func recvRes() error {
	idx, err := strconv.Atoi(models.Index)
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

	RecvNum = len(res.BlockList)

	idx, err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}
	RecvResList = append(RecvResList, *res)
	if len(res.BlockList) == len(RecvResList) {
		err = models.BlockChain[idx].SendRecvMd5Req(RecvResList)
		if err != nil {
			return nil
		}
		RecvResList = make([]models.RecoverRes, 0)
	}
	return nil
}

func ServeRecvMd5Req(data []byte) error {
	req := &models.RecoverMd5Req{}
	err := json.Unmarshal(data, req)
	if err != nil {
		return err
	}

	idx, err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}
	err = models.BlockChain[idx].SendRecvMd5Res(*req)
	if err != nil {
		return err
	}
	return nil
}

func ServeRecvMd5Res(data []byte) error {
	if RecvMd5ResList == nil {
		RecvMd5ResList = make([]models.RecoverMd5Res, 0)
	}
	res := &models.RecoverMd5Res{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return err
	}

	//idx,err := strconv.Atoi(models.Index)
	//if err != nil {
	//	return err
	//}
	RecvMd5ResList = append(RecvMd5ResList, *res)
	log.Println(RecvNum)
	log.Println(len(RecvMd5ResList))
	if RecvNum == len(RecvMd5ResList) {
		for i := 0; i < RecvNum; i++ {
			fmt.Println(RecvMd5ResList[i])
		}

		//err = models.BlockChain[idx].SendRecvMd5Req(RecvResList)
		//if err != nil {
		//	return nil
		//}
		RecvResList = make([]models.RecoverRes, 0)
	}
	return nil
}
