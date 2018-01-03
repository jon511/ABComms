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

func extractTagPacket(data []byte, programName string){
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

		tagLen := int(data[packetStart + 8])

		packet := data[packetStart:packetStart + tagLen + 10]

		//offsetArr := packet[0:2]
		//offsetArr = append(offsetArr, 0x00)
		//offsetArr = append(offsetArr, 0x00)

		//offset, err := bytesToInt32(true, offsetArr...)
		//if err != nil {
		//	fmt.Println(err)
		//}

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
	tag := LogixTag{}
	lenArr := packet[8:10]
	lenArr = append(lenArr, 0x00)
	lenArr = append(lenArr, 0x00)

	length, err := bytesToInt32(true, lenArr...)
	if err != nil {
		fmt.Println(err)
	}

	if programName != "" {
		tag.name = programName + "." + string(packet[10:length+10])
	}else{
		tag.name = string(packet[10:length + 10])
	}

	tagOffArr := packet[0:2]
	tagOffArr = append(tagOffArr, 0x00)
	tagOffArr = append(tagOffArr, 0x00)

	tag.offset, err = bytesToInt32(true , tagOffArr...)
	if err != nil {
		fmt.Println(err)
	}

	tag.dataType = string(packet[4])

	return tag
}
