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
	sequenceCounter int

	sessionHandle []byte
	otNetWorkConnectionID []byte
	responseData []byte

}


func (c Controller) connect() net.Conn{
	c.sequenceCounter = 1
	address := fmt.Sprintf("%s:%d", c.ipAddress, c.port)
	fmt.Println("connnecting to :", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
	}

	conn.Write(registerSession())

	conn.Read(c.responseData)
	copy(c.sessionHandle, c.responseData[4:8])

	conn.Write(forwardOpen(c.sessionHandle))
	conn.Read(c.responseData)

	copy(c.otNetWorkConnectionID, c.responseData[44:48])


	return conn
	//eipData := c.getTagList(conn)
	//
	//conn.Write(eipData)
	//myLen, _ := conn.Read(c.responseData)
	//
	//
	////sta := c.responseData[48:50]
	//status, _ := bytesToInt32(true, c.responseData[48], c.responseData[49], 0x00, 0x00)
	//
	//fmt.Println(status)
	//
	//
	//extractTagPacket(c.responseData[0:myLen], "")
	//
	//fmt.Println("")
	//fmt.Println(tagList)
	//fmt.Println(programNames)

}

func (c Controller) getTagList(conn net.Conn) []byte {

	//if !c.isConnected { return }

	_getTagList(c.sessionHandle, c.otNetWorkConnectionID, 1, conn)

	request := buildTagListRequest("",0)
	eipHeader := buildEipHeader(request, c.sessionHandle, c.otNetWorkConnectionID, 1)

	return eipHeader

}
