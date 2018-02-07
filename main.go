package main

import (

	"errors"
	"fmt"
	"bytes"
	"encoding/binary"
	"log"
)

var tagList []LogixTag
var programNames []string

var debug bool

func main() {

	//controller := Controller{"10.50.201.113", 44818, "", make([]byte,4)}
	controller := Controller{}
	controller.initializeController("10.50.193.55", 44818)
	//controller.initializeController("localhost", 44818)

	controller.connect()
	defer controller.conn.Close()

	//controller.conn.Close()

	//controller.testRead("rateFloat")
	tag := LogixTag{}
	tag.name = "rate"
	tag.controller = controller
	tag.dataType = DataType.DINT

	err := tag.read()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(tag.value)
	//fmt.Println(tag.dataType)
	//fmt.Println(tag.value)
	//fmt.Printf("%T", tag.value)

	tag.value = 200
	err = tag.write()
	if err != nil {
		log.Println(err)
	}


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
