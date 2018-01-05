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

	knownTags []string

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
	c.knownTags = make([]string, 0)
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


func (c Controller) testRead(tag string){

	requestService := byte(0x4c)

	requestPath := []byte{0x91, byte(len(tag))}
	requestPath = append(requestPath, []byte(tag)...)
	if len(tag) % 2 != 0 {
		requestPath = append(requestPath, 0x00)
	}
	requestPathSize := byte(len(requestPath)/2)
	requestData := []byte{0x01, 0x00}

	//sendData := []byte{0x4c, 0x03, 0x91, 0x04, 0x72, 0x61, 0x74, 0x65, 0x01, 0x00}
	var sendData []byte
	sendData = append(sendData, requestService)
	sendData = append(sendData, requestPathSize)
	sendData = append(sendData, requestPath...)
	sendData = append(sendData, requestData...)

	data := c.buildEipHeader(sendData)
	c.conn.Write(data)
	printHex(sendData)
	fmt.Println("data sent")

	resData := make([]byte, 1024)
	fmt.Println(resData)
	fmt.Println("reading data")
	len, err := c.conn.Read(resData)
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println(resData[0:len])
	printHex(resData[0:len])
}

func testWrite(c Controller, tag string, value int){

	requestService := byte(0x4D)

	requestPath := []byte{0x91, byte(len(tag))}
	requestPath = append(requestPath, []byte(tag)...)
	printHex(requestPath)
	requestPathSize := byte(len(requestPath) / 2)
	fmt.Println(requestPathSize)
	requestDataType := []byte{0xc4, 0x00}
	requestDataElements := []byte{0x01, 0x00}
	requestDataValue := int32ToSliceOfBytes(true, value, 4)



	tempSendData := []byte{0x4d, 0x06, 0x91, 0x0A, 0x43, 0x61, 0x72, 0x74, 0x6F, 0x6E, 0x53, 0x69, 0x7A, 0x65, 0xc4, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x00, 0x00}
	var sendData []byte
	sendData = append(sendData, requestService)
	sendData = append(sendData, requestPathSize)
	sendData = append(sendData, requestPath...)
	sendData = append(sendData, requestDataType...)
	sendData = append(sendData, requestDataElements...)
	sendData = append(sendData, requestDataValue...)

	printHex(sendData)
	printHex(tempSendData)
	data := c. buildEipHeader(sendData)
	c.conn.Write(data)

	fmt.Println("")
	fmt.Println(sendData)

	resData := make([]byte, 1024)

	fmt.Println("reading data")
	len, err := c.conn.Read(resData)
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println(resData[0:len])
	printHex(resData[0:len])
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

//func (c Controller) readTag(tag string){
//
//	//TODO: add parser to tag ffor use with arrays
//
//	t := tag
//	b := tag
//	i := 0
//
//
//
//}
//
//func TagNameParser(tag string, offset int){
//
//
//
//}
//
//func initialRead(tag, baseTag string, c Controller){
//
//	for _, val := range c.knownTags{
//		if val == baseTag {
//			return true
//		}
//	}
//
//	tagData := buildTagIOI(baseTag, false){
//
//	}
//
//}
//
//func buildTagIOI(tagName string, isBoolArray bool){
//	requestTagData := ""
//	tagArray := strings.Split(tagName, ".")
//
//	for i, _ := range tagArray{
//
//	}
//}
