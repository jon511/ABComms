package main

func registerSession() []byte {

	// from CIP Network Library Vol 2. Section 2-4.4.2.
	var data []byte
	//eip command
	data = append(data, []byte{0x65, 0x00}...)
	//eip length
	data = append(data, []byte{0x04, 0x00}...)
	//eip session handle
	data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)
	//eip status
	data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)
	//eip context
	data = append(data, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	//eip options
	data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)
	//eip protocol version
	data = append(data, []byte{0x01, 0x00}...)
	//eip option flay
	data = append(data, []byte{0x00, 0x00}...)

	return data
}

func forwardOpen(sessionHandle []byte) []byte {
	fwdOpen := cipForwardOpen()
	rrDataHeader := eipSendRRDataHeader(len(fwdOpen), sessionHandle)
	return append(rrDataHeader, fwdOpen...)
}

func cipForwardOpen() []byte {

	var data []byte
	// CIP Service
	data = append(data, 0x54)
	// CIP Path Size
	data = append(data, 0x02)
	// CIP Class Type
	data = append(data, 0x20)
	// CIP Class
	data = append(data, 0x06)
	// CIP Instance Type
	data = append(data, 0x24)
	// CIP Instance
	data = append(data, 0x01)
	// CIP Priority
	data = append(data, 0x0a)
	// CIP Timeout Ticks
	data = append(data, 0x0e)
	//CIP OT Connection ID (4 bytes)
	data = append(data, 0x02, 0x00, 0x00, 0x20)
	//CIP TO Connection ID (4 bytes)
	data = append(data, 0x01, 0x00, 0x00, 0x20)
	//CIP Connection Serial Number (2 bytes)
	data = append(data, int32ToSliceOfBytes(true, getRandomInt(63000), 2)...)
	//CIP Vendor ID (2 bytes)
	data = append(data, 0x37, 0x13)
	//DIP Originator serial number (4 bytes)
	data = append(data, 0x2a, 0x00, 0x00, 0x00)
	//cip multiplier (4 bytes)
	data = append(data, 0x03, 0x00, 0x00, 0x00)
	//cip OT RPI (4 bytes)
	data = append(data, 0x34, 0x12, 0x20, 0x00)
	//cip OT Network Connection Parameters (2 bytes)
	data = append(data, 0xf4, 0x43)
	//cip TO RPI (4 bytes)
	data = append(data, 0x01, 0x40, 0x20, 0x00)
	//cip TO Connection parameters (2 bytes)
	data = append(data, 0xf4, 0x43)
	//cip Transport trigger (1 byte)
	data = append(data, 0xa3)

	//connection path size / 2
	data = append(data, 0x03)

	//connection path
	data = append(data, 0x01, 0x00, 0x20, 0x02, 0x24, 0x01)

	return data

}

func eipSendRRDataHeader(frameLen int, sessionHandle []byte) []byte {

	var data []byte

	data = append(data,0x6f, 0x00)
	//eip length (2 bytes)
	data = append(data, int32ToSliceOfBytes(true, 16 + frameLen, 2)...)
	//data = append(data, 0x00)
	//eip session handle (4 bytes) value returned by RegisterSession
	data = append(data, sessionHandle...)
	//eip status (4 bytes)
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	//eip context (8 bytes)
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	//eip options (4 bytes)
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	//eip interface handle (4 bytes)
	data = append(data, 0x00, 0x00, 0x00, 0x00)
	//eip timeout (2 bytes)
	data = append(data, 0x00, 0x00)
	//eip item count (2 bytes)
	data = append(data, 0x02, 0x00)
	//eip item 1 type (2 bytes)
	data = append(data, 0x00, 0x00)
	//eip item 1 length (2 bytes)
	data = append(data, 0x00, 0x00)
	//eip item 2 type (2 bytes)
	data = append(data, 0xb2, 0x00)
	//eip item 2 length
	data = append(data, int32ToSliceOfBytes(true,frameLen, 2)...)

	return data
}

func buildTagListRequest(programName string, offset int) []byte {
	service := 0x55
	var pathSegment []byte

	if programName != "" {
		programNameBytes := []byte(programName)
		pathSegment = append(pathSegment, 0x91, byte(len(programName)))
		pathSegment = append(pathSegment, programNameBytes...)
		if len(programName) % 2 != 0 {
			pathSegment = append(pathSegment, 0x00)
		}

	}

	pathSegment = append(pathSegment, 0x20, 0x6b)

	if offset < 256 {
		pathSegment = append(pathSegment, 0x24)
	}else{
		pathSegment = append(pathSegment, 0x25)
	}

	pathSegmentLen := len(pathSegment) / 2
	attributeCount := []byte{0x03, 0x00}
	symbolType := []byte{0x02, 0x00}
	byteCount := []byte{0x07, 0x00}
	symbolName := []byte{0x01, 0x00}
	var attributes []byte //{ 0x03, 0x00, 0x02, 0x00, 0x07, 0x00, 0x01, 0x00}
	attributes = append(attributes, attributeCount...)
	attributes = append(attributes, symbolType...)
	attributes = append(attributes, byteCount...)
	attributes = append(attributes, symbolName...)

	tagListRequest := []byte{byte(service), byte(pathSegmentLen)}
	tagListRequest = append(tagListRequest, pathSegment...)
	tagListRequest = append(tagListRequest, attributes...)

	return tagListRequest

}

func buildEipHeader(tagIOI []byte, sessionHandle []byte) []byte {

	if contextHeader == 155 {
		contextHeader = 0
	}


	//eipPayLoadLength := byte(22 + len(tagIOI))
	eipConnectedDataLength := byte(len(tagIOI) + 2)

	eipCommand := []byte{0x70, 0x00}
	eipLength := []byte{byte(22 + len(tagIOI)), 0x00}
	eipSessionHandle := sessionHandle
	eipStatus := []byte {0x00, 0x00, 0x00, 0x00}
	//TODO: need to figure out what this is
	eipContext := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	eipOptions := []byte{0x00, 0x00, 0x00, 0x00}
	eipInterfaceHandle := []byte{0x00, 0x00, 0x00, 0x00}
	eipTimeout := []byte{0x00, 0x00}
	eipItemCount := []byte{0x02, 0x00}
	eipItem1ID := []byte{0xa1, 0x00}
	eipItem1Length := []byte{0x04, 0x00}
	//TODO: need to figure this out
	eipItem1 := []byte{} //OTNetworkConnectionID
	eipItem2ID := []byte{0xb1, 0x00}
	eipItem2Length := []byte{byte(eipConnectedDataLength), 0x00}
	//TODO: need to find out about the sequence counter and it's purpose.
	eipSequence := []byte{0x00, 0x00} //sequence counter

	var eipHeaderFrame []byte
	eipHeaderFrame = append(eipHeaderFrame, eipCommand...)
	eipHeaderFrame = append(eipHeaderFrame, eipLength...)
	eipHeaderFrame = append(eipHeaderFrame, eipSessionHandle...)
	eipHeaderFrame = append(eipHeaderFrame, eipStatus...)
	eipHeaderFrame = append(eipHeaderFrame, eipContext...)
	eipHeaderFrame = append(eipHeaderFrame, eipOptions...)
	eipHeaderFrame = append(eipHeaderFrame, eipInterfaceHandle...)
	eipHeaderFrame = append(eipHeaderFrame, eipTimeout...)
	eipHeaderFrame = append(eipHeaderFrame, eipItemCount...)
	eipHeaderFrame = append(eipHeaderFrame, eipItem1ID...)
	eipHeaderFrame = append(eipHeaderFrame, eipItem1Length...)
	eipHeaderFrame = append(eipHeaderFrame, eipItem1...)
	eipHeaderFrame = append(eipHeaderFrame, eipItem2ID...)
	eipHeaderFrame = append(eipHeaderFrame, eipItem2Length...)
	eipHeaderFrame = append(eipHeaderFrame, eipSequence...)

	return append(eipHeaderFrame, tagIOI...)

}

var context_dict = map[int] []byte{
	0: {0x57, 0x65, 0x27, 0x72, 0x65, 0x00, 0x00, 0x00},
	1: {0x6e, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	2: {0x73, 0x74, 0x72, 0x61, 0x6e, 0x67, 0x00, 0x00},
	3: {0x65, 0x72, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00},
	4: {0x74, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	5: {0x6c, 0x6f, 0x76, 0x65, 0x00, 0x00, 0x00, 0x00},
	6: {0x59, 0x6f, 0x75, 0x00, 0x00, 0x00, 0x00, 0x00},
	7: {0x6b, 0x6e, 0x6f, 0x77, 0x00, 0x00, 0x00, 0x00},
	8: {0x74, 0x68, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00},
	9: {0x72, 0x75, 0x6c, 0x65, 0x73, 0x00, 0x00, 0x00},
	10: {0x61, 0x6e, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00},
	11: {0x73, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	12: {0x64, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	13: {0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	14: {0x41, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	15: {0x66, 0x75, 0x6c, 0x6c, 0x00, 0x00, 0x00, 0x00},
	16: {0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x00, 0x00},
	17: {0x6d, 0x65, 0x6e, 0x74, 0x27, 0x73, 0x00, 0x00},
	18: {0x77, 0x68, 0x61, 0x74, 0x00, 0x00, 0x00, 0x00},
	19: {0x49, 0x27, 0x6d, 0x00, 0x00, 0x00, 0x00, 0x00},
	20: {0x74, 0x68, 0x69, 0x6e, 0x6b, 0x00, 0x00, 0x00},
	21: {0x69, 0x6e, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00},
	22: {0x6f, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	23: {0x59, 0x6f, 0x75, 0x00, 0x00, 0x00, 0x00, 0x00},
	24: {0x77, 0x6f, 0x75, 0x6c, 0x64, 0x6e, 0x74, 0x00},
	25: {0x67, 0x65, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00},
26: 0x73696874,
27: 0x6d6f7266,
28: 0x796e61,
29: 0x726568746f,
30: 0x797567,
31: 0x49,
32: 0x7473756a,
33: 0x616e6e6177,
34: 0x6c6c6574,
35: 0x756f79,
36: 0x776f68,
37: 0x6d2749,
38: 0x676e696c656566,
39: 0x6174746f47,
40: 0x656b616d,
41: 0x756f79,
42: 0x7265646e75,
43: 0x646e617473,
44: 0x726576654e,
45: 0x616e6e6f67,
46: 0x65766967,
47: 0x756f79,
48: 0x7075,
49: 0x726576654e,
50: 0x616e6e6f67,
51: 0x74656c,
52: 0x756f79,
53: 0x6e776f64,
54: 0x726576654e,
55: 0x616e6e6f67,
56: 0x6e7572,
57: 0x646e756f7261,
58: 0x646e61,
59: 0x747265736564,
60: 0x756f79,
61: 0x726576654e,
62: 0x616e6e6f67,
63: 0x656b616d,
64: 0x756f79,
65: 0x797263,
66: 0x726576654e,
67: 0x616e6e6f67,
68: 0x796173,
69: 0x657962646f6f67,
70: 0x726576654e,
71: 0x616e6e6f67,
72: 0x6c6c6574,
73: 0x61,
74: 0x65696c,
75: 0x646e61,
76: 0x74727568,
77: 0x756f79,
78: 0x6576276557,
79: 0x6e776f6e6b,
80: 0x68636165,
81: 0x726568746f,
82: 0x726f66,
83: 0x6f73,
84: 0x676e6f6c,
85: 0x72756f59,
86: 0x73277472616568,
87: 0x6e656562,
88: 0x676e69686361,
89: 0x747562,
90: 0x657227756f59,
91: 0x6f6f74,
92: 0x796873,
93: 0x6f74,
94: 0x796173,
95: 0x7469,
96: 0x656469736e49,
97: 0x6577,
98: 0x68746f62,
99: 0x776f6e6b,
100: 0x732774616877,
101: 0x6e656562,
102: 0x676e696f67,
103: 0x6e6f,
104: 0x6557,
105: 0x776f6e6b,
106: 0x656874,
107: 0x656d6167,
108: 0x646e61,
109: 0x6572276577,
110: 0x616e6e6f67,
111: 0x79616c70,
112: 0x7469,
113: 0x646e41,
114: 0x6669,
115: 0x756f79,
116: 0x6b7361,
117: 0x656d,
118: 0x776f68,
119: 0x6d2749,
120: 0x676e696c656566,
121: 0x74276e6f44,
122: 0x6c6c6574,
123: 0x656d,
124: 0x657227756f79,
125: 0x6f6f74,
126: 0x646e696c62,
127: 0x6f74,
128: 0x656573,
129: 0x726576654e,
130: 0x616e6e6f67,
131: 0x65766967,
132: 0x756f79,
133: 0x7075,
134: 0x726576654e,
135: 0x616e6e6f67,
136: 0x74656c,
137: 0x756f79,
138: 0x6e776f64,
139: 0x726576654e,
140: 0x6e7572,
141: 0x646e756f7261,
142: 0x646e61,
143: 0x747265736564,
144: 0x756f79,
145: 0x726576654e,
146: 0x616e6e6f67,
147: 0x656b616d,
148: 0x756f79,
149: 0x797263,
150: 0x726576654e,
151: 0x616e6e6f67,
152: 0x796173,
153: 0x657962646f6f67,
154: 0x726576654e,
155: 0xa680e2616e6e6f67}