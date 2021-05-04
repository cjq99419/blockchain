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
	"sync"
	"time"
)

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PreHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func addBlock(idx int, data string, port string, lon string, lat string) (Block, error) {
	var newBlock Block
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data

	newBlock.Addr = Address{
		Ip:   "127.0.0.1",
		Port: port,
	}

	newBlock.Token = "safe token"

	newBlock.Type = 0

	if BlockChain == nil || len(BlockChain) == 0 {
		BlockChain = make([]Block, 0)
		newBlock.Index = idx
		newBlock.PreHash = ""
		newBlock.Lon = "9999"
		newBlock.Lat = "9999"
		newBlock.Hash = calculateHash(newBlock)
		BlockChain = append(BlockChain, newBlock)
		newBlock.Data = ""

		_, err := addBlock(idx, data, port, lon, lat)
		if err != nil {
			return Block{}, err
		}
		BlockChain = BlockChain[1:]
	} else {
		oldBlock := BlockChain[len(BlockChain)-1]
		newBlock.Id = strconv.Itoa(oldBlock.Index + 1)
		newBlock.Index = idx
		newBlock.Lon = lon
		newBlock.Lat = lat
		newBlock.PreHash = oldBlock.Hash
		newBlock.Hash = calculateHash(newBlock)
		BlockChain = append(BlockChain, newBlock)
	}
	return newBlock, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Println("[Error]:port is not found")
	}
	Index = os.Args[1]
	HTTPPort = fmt.Sprintf("800%v", Index)
	GRPCPort = fmt.Sprintf("801%v", Index)

	var err error
	_, err = addBlock(0, "aaaaaaafasdfasdfasdf", "8000", "1", "2")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(1, "aaaaaaafasdfasdfasdf", "8001", "2", "6")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(2, "aaaaaaafasdfasdfasdf", "8002", "7", "3")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(3, "aaaaaaafasdfasdfasdf", "8003", "10", "2")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(4, "aaaaaaafasdfasdfasdf", "8004", "15", "2")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(5, "aaaaaaafasdfasdfasdf", "8005", "1", "21")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(6, "aaaaaaafasdfasdfasdf", "8006", "10", "22")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(7, "aaaaaaafasdfasdfasdf", "8007", "18", "21")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(8, "aaaaaaafasdfasdfasdf", "8008", "4", "9")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	_, err = addBlock(9, "aaaaaaafasdfasdfasdf", "8009", "12", "22")
	if err != nil {
		log.Printf("[Error]:%v", err)
	}
	for _, e := range BlockChain {
		fmt.Println(e)
	}

	err = InitDist()
	if err != nil {
		log.Printf("[Error]:%v", err)
	} else {
		log.Printf("[Info]:init dist successful\n")
	}

	//idx, _ := strconv.Atoi(Index)
	//fmt.Println(Index)
	//BlockChain[idx].Type = 2
	//err = BlockChain[idx].SendRecvReq()
	//if err != nil {
	//	log.Printf("[Error]:%v", err)
	//} else {
	//	log.Printf("[Info]:OK\n")
	//}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = service.StartHttpService(HTTPPort)
		if err != nil {
			log.Printf("[Error]:%v", err)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = service.StartGRPCService(GRPCPort)
		if err != nil {
			log.Printf("[Error]:%v", err)
		}
	}(&wg)

	wg.Wait()

}
