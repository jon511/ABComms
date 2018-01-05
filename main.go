package main

import (

	"errors"
	"math/rand"
	"time"
	"fmt"
	"math"
	"bytes"
	"encoding/binary"
)

var tagList []LogixTag
var programNames []string

func main() {

	fmt.Println(math.Float32frombits(1084227584))
	i := math.Float32bits(5.0)
	fmt.Println(i)

	in := read_int32([]byte{0xff, 0x0f, 0x00, 0x00})
	fmt.Println(in)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(in))
	printHex(bs)


	//controller := Controller{"10.50.201.113", 44818, "", make([]byte,4)}
	controller := Controller{}
	controller.initializeController("10.50.193.55", 44818)
	fmt.Println(controller.sequenceCounter)
	controller.connect()

	//controller.conn.Close()

	//controller.testRead("rateFloat")
	tag := LogixTag{}
	tag.name = "rateFloat"
	tag.controller = controller

	tag.read()
	fmt.Println(tag.dataType)
	fmt.Println(tag.value)
	fmt.Printf("%T", tag.value)

	tag.value = 1896.56
	tag.write()

	//testWrite(controller, "testDint", 2500)



	//controller.getTagList()
	//fmt.Println(tagList)
	//fmt.Println(programNames)

}

func read_int32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
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