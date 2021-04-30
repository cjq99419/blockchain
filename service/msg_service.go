package service

import (
	"blockchain/models"
	"encoding/json"
	"log"
	"strconv"
)

func RecvVerifyReq(data []byte) error {
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
	err = models.BlockChain[idx].RecvVerifyReq(*req)
	if err != nil {
		return err
	}
	return nil
}