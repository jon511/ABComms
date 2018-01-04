package main

import (
	"fmt"
	"net"
	"log"
	"errors"
)

type Controller struct{

	ipAddress string
	processorSlot int
	micro800 bool
	port int
	vendorID []byte
	context []byte
	contextPointer int

	conn net.Conn
	setTimeout float64
	isConnected bool
	otNetWorkConnectionID []byte
	sessionHandle []byte
	sessionIsRegistered bool
	serialNumber []byte
	originatorSerialNumber []byte
	sequenceCounter int
	offset int
	structIdentifier []byte

	initialized bool

}

func (c *Controller) initializeController(ipAddress string, port int){
	c.ipAddress = ipAddress
	c.port = port
	c.vendorID = []byte{0x37, 0x13}
	c.sessionHandle = []byte{0x00, 0x00, 0x00, 0x00}
	c.serialNumber = int32ToSliceOfBytes(true, getRandomInt(63000), 2)
	c.context = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	c.otNetWorkConnectionID = []byte{0x00, 0x00, 0x00, 0x00}
	c.originatorSerialNumber = []byte{0x42, 0x00, 0x00, 0x00}
	c.structIdentifier = []byte{0xce, 0x0f}
	c.sequenceCounter = 1
	c.initialized = true
}

func (c *Controller) connect() error {

	if c.ipAddress == "" {
		return errors.New("ip address of controller cannot be nil")
	}

	if c.port == 0{
		c.port = 44818
	}

	address := fmt.Sprintf("%s:%d", c.ipAddress, c.port)
	fmt.Println("connnecting to :", address)


	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
	}

	c.conn = conn

	c.conn.Write(c.buildRegisterSession())

	var responseData = make([]byte, 1028)

	_, err = c.conn.Read(responseData)
	if err != nil {
		log.Println(err)
		c.isConnected = false
	}else{
		copy(c.sessionHandle, responseData[4:8])
	}
	printHex(c.sessionHandle)
	printHex(c.forwardOpenPacket())
	c.conn.Write(c.forwardOpenPacket())
	_, err = c.conn.Read(responseData)
	if err != nil {
		log.Println("Foward open failed")
		c.isConnected = false
	}else{
		copy(c.otNetWorkConnectionID, responseData[44:48])
		c.isConnected = true
	}
	printHex(responseData)

	return nil

}

//func (c Controller) getTagList() []byte {
//
//	//if !c.isConnected { return }
//
//	_getTagList(c.sessionHandle, c.otNetWorkConnectionID, 1, conn)
//
//	request := buildTagListRequest("",0)
//	eipHeader := buildEipHeader(request, c.sessionHandle, c.otNetWorkConnectionID, 1)
//
//	return eipHeader
//
//}
