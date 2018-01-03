package main

import (
	"fmt"
	"net"
	"log"
)

type Controller struct{
	ipAddress string
	port int

	isConnected bool

	sessionHandle []byte
	responseData []byte

}


func (c Controller) connect(){
	address := fmt.Sprintf("%s:%d", c.ipAddress, c.port)
	fmt.Println("connnecting to :", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
	}

	printHex(registerSession())

	conn.Write(registerSession())

	conn.Read(c.responseData)
	copy(c.sessionHandle, c.responseData[4:8])

	printHex(c.sessionHandle)

	printHex(c.responseData)

	printHex(forwardOpen(c.sessionHandle))

	conn.Write(forwardOpen(c.sessionHandle))
	conn.Read(c.responseData)

	printHex(c.responseData)

	eipData := c.getTagList()

	printHex(eipData)

	conn.Write(eipData)
	conn.Read(c.responseData)

	printHex(c.responseData)

}

func (c Controller) getTagList() []byte {

	//if !c.isConnected { return }

	programNames := make(map[string]string)
	tagList := make(map[string]string)

	programNames["program"] = "name"
	tagList["tag"] = "list"



	request := buildTagListRequest("",0)
	eipHeader := buildEipHeader(request, c.sessionHandle)

	return eipHeader

}
