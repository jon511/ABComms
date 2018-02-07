package main

import (
	"fmt"
	"strings"
	"encoding/binary"
	"math"
	"bytes"
	"errors"
)

type DataTypeStruct struct {
	name string
	value int
}

type LogixDataTypes struct {
	BOOL DataTypeStruct
	SINT DataTypeStruct
	INT DataTypeStruct
	DINT DataTypeStruct
	REAL DataTypeStruct
	DWORD DataTypeStruct
	LINT DataTypeStruct
}

var DataType = LogixDataTypes {
	DataTypeStruct{"BOOL",0xc1},
	DataTypeStruct{"SINT",0xc2},
	DataTypeStruct{"INT",0xc3},
	DataTypeStruct{"DINT",0xc4},
	DataTypeStruct{"REAL", 0xca},
	DataTypeStruct{"DWORD",0xd3},
	DataTypeStruct{"LINT", 0xc5},
	}

type LogixTag struct{
	name string
	offset int
	dataType DataTypeStruct
	controller Controller
	qualityCode int
	value interface{}

}

func (t *LogixTag) read() error{

	if !t.controller.isValid() {
		return errors.New("tag controller field is not valid")
	}

	if !t.controller.isConnected {
		return errors.New("tag controller is not connected")
	}

	//TODO add error checking before attempting to read tag
	requestService := byte(0x4c)
	requestPath := []byte{0x91, byte(len(t.name))}
	requestPath = append(requestPath, []byte(t.name)...)
	if len(t.name) % 2 != 0 {
		requestPath = append(requestPath, 0x00)
	}
	requestPathSize := byte(len(requestPath)/2)
	requestData := []byte{0x01, 0x00}

	var sendData []byte
	sendData = append(sendData, requestService)
	sendData = append(sendData, requestPathSize)
	sendData = append(sendData, requestPath...)
	sendData = append(sendData, requestData...)
	fmt.Println("send data")
	fmt.Println(sendData)
	data := t.controller.buildEipHeader(sendData)

	printHex(sendData)
	fmt.Println(len(data))
	t.controller.conn.Write(data)
	printHex(data)
	resData := make([]byte, 1024)
	retLen, err := t.controller.conn.Read(resData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("read data")
	printHex(data)

	startPointer := 0

	for i, b := range resData{
		if b == 0xcc{
			startPointer = i
			break
		}
	}

	response := resData[startPointer:retLen]

	//TODO add tag verification check on read

	switch response[4]{
	case 0xc1:
		t.dataType = DataType.BOOL
	case 0xc2:
		t.dataType = DataType.SINT
	case 0xc3:
		t.dataType = DataType.INT
	case 0xc4:
		t.dataType = DataType.DINT
		var retVal int32
		buf := bytes.NewBuffer(response[6:])
		binary.Read(buf, binary.LittleEndian, &retVal)
		t.value = retVal
	case 0xca:
		t.dataType = DataType.REAL
		var retVal uint32
		buf := bytes.NewBuffer(response[6:])
		binary.Read(buf, binary.LittleEndian, &retVal)
		t.value = math.Float32frombits(retVal)
	case 0xd3:
		t.dataType = DataType.DWORD
	case 0xc5:
		t.dataType = DataType.LINT

	}

	printHex(response)

	return nil
}

func (t LogixTag) write() error{

	if !t.controller.isConnected {
		return errors.New("tag controller is not connected")
	}

	fmt.Printf("%T", t.value)

	//TODO add error checking before attempting to write tag

	writeValue := make([]byte, 4)

	switch t.dataType.name{
	case "BOOL":
		t.dataType = DataType.BOOL
	case "SINT":
		t.dataType = DataType.SINT
	case "INT":
		t.dataType = DataType.INT
	case "DINT":
		binary.LittleEndian.PutUint32(writeValue, uint32(t.value.(int)))
	case "REAL":
		f := float32(t.value.(float64))
		i := math.Float32bits(f)
		binary.LittleEndian.PutUint32(writeValue, i)
	case "DWORD":
		t.dataType = DataType.DWORD
	case "LINT":
		t.dataType = DataType.LINT

	}

	requestService := byte(0x4D)

	requestPath := []byte{0x91, byte(len(t.name))}
	requestPath = append(requestPath, []byte(t.name)...)
	if len(requestPath) % 2 != 0 {
		requestPath = append(requestPath, 0x00)
	}
	requestPathSize := byte(len(requestPath) / 2)

	requestDataType := []byte{byte(t.dataType.value), 0x00}
	requestDataElements := []byte{0x01, 0x00}
	requestDataValue := writeValue

	var sendData []byte
	sendData = append(sendData, requestService)
	sendData = append(sendData, requestPathSize)
	sendData = append(sendData, requestPath...)
	sendData = append(sendData, requestDataType...)
	sendData = append(sendData, requestDataElements...)
	sendData = append(sendData, requestDataValue...)


	data := t.controller.buildEipHeader(sendData)
	t.controller.conn.Write(data)

	fmt.Println("write data")
	printHex(data)

	resData := make([]byte, 1024)

	_, err := t.controller.conn.Read(resData)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(resData[0:len])
	//printHex(resData[0:len])

	return nil
}

func (c *Controller) extractTagPacket(data []byte, programName string){
	packetStart := 50

	for packetStart < len(data) {

		//tagLenArr := data[packetStart + 8: packetStart + 10]

		//
		//tagLenArr = append(tagLenArr, 0x00)
		//tagLenArr = append(tagLenArr, 0x00)
		//
		//tagLen, err := bytesToInt32(true, tagLenArr...)
		//if err != nil {
		//	fmt.Println(err)
		//}

		//tagLen := int(data[packetStart + 8])

		tagLen, err := bytesToInt32(true, data[packetStart + 8], data[packetStart + 9], 0x00, 0x00)
		if err != nil {
			fmt.Println(err)
		}

		packet := data[packetStart:packetStart + tagLen + 10]

		//offset, err := bytesToInt32(true, packet[0], packet[1], 0x00, 0x00)
		//if err != nil {
		//	fmt.Println(err)
		//}

		//c.offset = offset
		//fmt.Println(offset)

		tag := parseLgxTag(packet, "")

		if !strings.Contains(tag.name, "__DEFVAL_") && !strings.Contains(tag.name, "Routine:"){
			tagList = append(tagList, tag)
		}

		if programName == "" {
			if strings.Contains(tag.name, "Program:"){
				programNames = append(programNames, tag.name)
			}
		}

		packetStart = packetStart + tagLen + 10

	}
}

func parseLgxTag(packet []byte, programName string) LogixTag{
	fmt.Println(string(packet))
	fmt.Println(packet)
	tag := LogixTag{}

	length, err := bytesToInt32(true, packet[8], packet[9], 0x00, 0x00)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(length)
	if programName != "" {
		tag.name = programName + "." + string(packet[10:length+10])
	}else{
		tag.name = string(packet[10:length + 10])
	}

	//tagOffArr := packet[0:2]
	//tagOffArr = append(tagOffArr, 0x00)
	//tagOffArr = append(tagOffArr, 0x00)
	//
	//tag.offset, err = bytesToInt32(true , tagOffArr...)
	//if err != nil {
	//	fmt.Println(err)
	//}

	tag.offset = int(packet[0])

	//tag.dataType = string(packet[4])
	fmt.Println(tag.name)
	return tag
}
