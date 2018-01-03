package main

import (

	"errors"
	"math/rand"
	"time"
	"fmt"
)

var contextHeader int


func main() {


	//controller := Controller{"10.50.201.113", 44818, "", make([]byte,4)}
	controller := Controller{}
	controller.ipAddress = "10.50.201.113"
	controller.port = 44818
	controller.sessionHandle = make([]byte, 4)
	controller.responseData = make([] byte, 1028)

	controller.connect()
	controller.getTagList()

}

func printHex(data []byte){
	fmt.Printf("[ ")
	for _, i := range data {
		fmt.Printf("0x%x ", i)
	}
	fmt.Printf("]\n")

}

func int32ToSliceOfBytes(littleEndian bool, i, minimumBytesToReturn int) []byte{

	var b []byte

	switch {
	case i < 0x000000ff:
		b = append(b, byte(i))
	case i > 0x000000ff && i < 0x00010000:
		b = append(b, byte(i & 0x000000ff))
		b = append(b, byte((i & 0x0000ff00) >> 8))
	case i > 0x0000ff00 && i < 0x01000000:
		b = append(b, byte(i & 0x000000ff))
		b = append(b, byte((i & 0x0000ff00) >> 8))
		b = append(b, byte((i & 0x00ff0000) >> 16))
	case i > 0x00ff0000:
		b = append(b, byte(i & 0x000000ff))
		b = append(b, byte((i & 0x0000ff00) >> 8))
		b = append(b, byte((i & 0x00ff0000) >> 16))
		b = append(b, byte((i & 0x7f000000) >> 24))

	}
	for i := len(b); i < minimumBytesToReturn; i++ {
		b = append(b, 0x00)
	}

	return b
}

func bytesToInt32(littleEndian bool, b ...byte) (int, error) {

	if len(b) != 4 {
		return -1, errors.New("bytesToInt32 : slice of bytes must contain 4 elements")
	}

	if littleEndian {
		return int(b[3]) << 24 + int(b[2]) << 16 + int(b[1]) << 8 + int(b[0]), nil
	}

	return int(b[0]) << 24 + int(b[1]) << 16 + int(b[2] << 8) + int(b[3]), nil

}

func getRandomInt(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max)
}