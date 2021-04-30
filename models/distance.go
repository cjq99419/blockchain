package models

import (
	"errors"
	"math"
	"strconv"
)

var Dist [][]float32

func InitDist() error {
	lens := len(BlockChain)
	var err error
	Dist = make([][]float32, lens)
	for i := 0; i < lens; i++ {
		Dist[i] = make([]float32, lens)
		for j := 0; j < lens; j++ {
			Dist[i][j], err = calculate(i, j, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateDist(idx int) error {
	for i := 0; i < len(BlockChain); i++ {
		dist, err := calculate(i, idx, true)
		if err != nil {
			return err
		}
		Dist[i][idx] = dist
		Dist[idx][i] = dist
	}
	return nil
}

func calculate(idx1, idx2 int, redo bool) (float32, error) {
	if Dist[idx2]!= nil && Dist[idx2][idx1] != 0 && !redo {
		return Dist[idx2][idx1], nil
	}
	if idx2 == idx1 {
		return 0, nil
	}

	block1 := BlockChain[idx1]
	block2 := BlockChain[idx2]

	lon1, err := strconv.ParseFloat(block1.Lon, 32)
	if err != nil {
		return 0, err
	}
	lat1, err := strconv.ParseFloat(block1.Lat, 32)
	if err != nil {
		return 0, err
	}
	lon2, err := strconv.ParseFloat(block2.Lon, 32)
	if err != nil {
		return 0, err
	}
	lat2, err := strconv.ParseFloat(block2.Lon, 32)
	if err != nil {
		return 0, err
	}

	return float32(math.Sqrt(math.Pow(lon1-lon2, 2) + math.Pow(lat1-lat2, 2))), nil
}

func GetTopNDist(idx int, N int) ([]int, error) {
	if N >= len(BlockChain) {
		return nil, errors.New("failed to get top n dist")
	}

	distList := make([]int, len(BlockChain))
	nodeList := make([]int, len(BlockChain))

	result := make([]int, N)

	for i := 0; i < len(BlockChain); i++ {
		nodeList[i] = i
		distList[i] = int(Dist[i][idx])
	}
	quickSort(distList, nodeList)

	for i := 0; i < N; i++ {
		result[i] = nodeList[i+1]
	}

	return result, nil
}

func quickSort(distList []int, nodeList []int) {
	if len(distList) <= 1 {
		return
	}
	flag := distList[0]
	left, right := 0, len(distList)-1

	for i := 1; i <= right; {
		if distList[i] > flag {
			distList[i], distList[right] = distList[right], distList[i]
			nodeList[i], nodeList[right] = nodeList[right], nodeList[i]
			right--
		} else {
			distList[i], distList[left] = distList[left], distList[i]
			nodeList[i], nodeList[left] = nodeList[left], nodeList[i]
			i++
			left++
		}
	}
	// 递归
	quickSort(distList[:left], nodeList[:left])
	quickSort(distList[left+1:], nodeList[left+1:])
}
