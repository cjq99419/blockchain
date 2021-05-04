package service

import (
	"blockchain/models"
	"blockchain/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"strconv"
	"time"
)

var DataArr map[int64][]byte

type DownloadServer struct{}

func StartGRPCService(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	log.Printf("[Info]:Start grpc listen to 127.0.0.1:%v", port)
	//构建一个新的服务端对象
	s := grpc.NewServer()
	//向这个服务端对象注册服务
	proto.RegisterDownloadServer(s, &DownloadServer{})
	//注册服务端反射服务
	reflection.Register(s)

	//启动服务
	return s.Serve(lis)
}

func (*DownloadServer) Download(req *proto.DownloadReq, downloadServer proto.Download_DownloadServer) error {
	offset := req.Offset
	blockSize := int64(64 * 1024)
	//循环发送数据
	for {
		if offset > req.Offset+req.Size {
			break
		} else if offset+int64(blockSize) >= req.Offset+req.Size {
			res := &proto.DownloadRes{
				Offset: offset,
				Size:   req.Offset + req.Size - offset,
			}
			res.Data = []byte(models.BlockChain[req.Base.To].Data)[res.Offset : res.Offset+res.Size]
			err := downloadServer.Send(res)
			if err != nil {
				return err
			}
			break
		} else {
			err := downloadServer.Send(&proto.DownloadRes{
				Offset: offset,
				Size:   blockSize,
				Data:   []byte(models.BlockChain[req.Base.To].Data)[offset : offset+blockSize],
			})
			if err != nil {
				return err
			}
			offset += blockSize
		}
	}
	return nil
}

func GRPCDataService(reqs []models.DownloadReq) error {
	log.Printf("[Info]:Start grpc service")
	DataArr = make(map[int64][]byte)
	//wg := sync.WaitGroup{}
	for _, req := range reqs {
		//go func(wg *sync.WaitGroup) {
		//	wg.Add(1)
		//	defer wg.Done()
		err := gRPCDownload(req)
		if err != nil {
			log.Printf("[Error]:%v", err)
		}
		//	}(&wg)
		//}
		//wg.Wait()
	}
	return dataRecovery()
}

func dataRecovery() error {
	keys := make([]int,0)
	for k := range DataArr {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	data := make([]byte,0)
	for _, key := range keys {
		data = append(data, DataArr[int64(key)]...)
	}

	idx,err := strconv.Atoi(models.Index)
	if err != nil {
		return err
	}
	models.BlockChain[idx].Data = string(data)
	log.Printf("[Info]:Recovery successful!")
	return nil
}

func gRPCDownload(req models.DownloadReq) error {
	grpcConn, err := grpc.Dial(models.BlockChain[req.To].Addr.Ip+":"+fmt.Sprintf("801%v",models.BlockChain[req.To].Addr.Port[len(models.BlockChain[req.To].Addr.Port) - 1:]), grpc.WithInsecure())
	if err != nil {
		return err
	}

	//通过grpc连接创建一个客户端实例对象
	client := proto.NewDownloadClient(grpcConn)

	//设置ctx超时（根据情况设定）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//和简单rpc不同，此时获得的不是res，而是一个client的对象，通过这个连接对象去读取数据
	downloadClient, err := client.Download(ctx, &proto.DownloadReq{
		Base: &proto.BaseMessage{
			MessageId:   req.MessageId,
			From:        int32(req.From),
			To:          int32(req.To),
			MessageName: req.MessageName,
			Timestamp:   req.Timestamp,
		},
		Offset: req.Slice.Offset,
		Size:   req.Slice.Size,
	})
	if err != nil {
		return err
	}
	data := make([]byte, req.Slice.Size)
	var sumSize int64
	//循环处理数据，当监测到读取完成后退出
	for {
		res, err := downloadClient.Recv()
		if err != nil {
			return err
		}
		sumSize += res.Size
		data = append(data, res.Data...)
		log.Printf("[Info]:Get a date package, offset:%v, size:%v\n", res.Offset, res.Size)
		if sumSize >= req.Slice.Size {
			break
		}
	}
	DataArr[req.Slice.Offset] = data
	return nil
}
