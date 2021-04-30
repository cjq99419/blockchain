package main

import (
	. "blockchain/models"
	"blockchain/service"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PreHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func addBlock(data string, port string) (Block, error) {
	var newBlock Block
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.Lon = "0"
	newBlock.Lat = "0"

	newBlock.Addr = Address{
		Ip:   "127.0.0.1",
		Port: port,
	}

	newBlock.Token = "safe token"

	newBlock.Type = 0


	if BlockChain == nil || len(BlockChain) == 0 {
		BlockChain = make([]Block, 0)
		newBlock.Index = 0
		newBlock.PreHash = ""
		newBlock.Hash = calculateHash(newBlock)
		BlockChain = append(BlockChain, newBlock)
		newBlock.Data = ""

		_, err := addBlock(data,port)
		if err != nil {
			return Block{},err
		}
	} else {
		oldBlock := BlockChain[len(BlockChain)-1]
		newBlock.Id = strconv.Itoa(oldBlock.Index + 1)
		newBlock.Index = oldBlock.Index + 1


		newBlock.PreHash = oldBlock.Hash
		newBlock.Hash = calculateHash(newBlock)
		BlockChain = append(BlockChain, newBlock)
	}
	return newBlock, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Println("[error]:port is not found")
	}
	Index = os.Args[1]
	Port = fmt.Sprintf("800%v",Index)

	var err error



	_, err = addBlock("aaaaaaafasdfasdfasdf", "8001")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	_, err = addBlock("aaaaaaafasdfasdfasdf", "8002")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	_, err = addBlock("aaaaaaafasdfasdfasdf", "8003")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	_, err = addBlock("aaaaaaafasdfasdfasdf", "8004")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	_, err = addBlock("aaaaaaafasdfasdfasdf", "8005")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	_, err = addBlock("aaaaaaafasdfasdfasdf", "8006")
	if err != nil {
		log.Printf("[error]:%v",err)
	}
	for _, e := range BlockChain {
		fmt.Println(e)
	}

	err = InitDist()
	if err != nil {
		log.Printf("[Error]:%v",err)
	} else {
		log.Printf("[Info]:init dist successful\n")
	}

	idx,_ := strconv.Atoi(Index)
	fmt.Println(Index)
	BlockChain[idx].Type = 2
	err = BlockChain[idx].SendRecvReq()
	if err != nil {
		log.Printf("[Error]:%v",err)
	} else {
		log.Printf("[Info]:OK\n")
	}

	service.StartHttpService(Port)


}
