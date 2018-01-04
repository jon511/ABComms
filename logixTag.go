package main

import (
	"fmt"
	"strings"
)

type LogixTag struct{
	name string
	offset int
	dataType string
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

	tag.dataType = string(packet[4])
	fmt.Println(tag.name)
	return tag
}
