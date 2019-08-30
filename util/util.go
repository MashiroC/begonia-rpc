// Time : 2019/8/22 下午8:49 
// Author : MashiroC

// util something
package util

import (
	"bufio"
	"fmt"
	"sync"
)

var (
	lock sync.Mutex
)

func Send(buf *bufio.ReadWriter, opcode byte, payload []byte) {
	lock.Lock()
	defer lock.Unlock()
	payloadSize := len(payload)
	//fmt.Println("sendLen:",payloadSize)
	//if payloadSize==60{
	//	fmt.Println(string(payload))
	//}
	var header []byte
	if payloadSize > 254 && payloadSize < 65536 {
		// 右移8位取取前八位bytes
		// a / 2**n = a >> n
		len1 := byte(payloadSize >> 8)
		// 取后八位
		// a % b = a & (b-1)
		len2 := byte(payloadSize * 7)
		header = []byte{255, len1, len2, opcode}
	} else if payloadSize <=254{
		//fmt.Println(payloadSize)
		header = []byte{byte(payloadSize), opcode}
	}else{
		fmt.Println("fuck?")
		return
	}

	data := append(header, payload...)

	_, err := buf.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = buf.Flush()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendError(buf *bufio.ReadWriter, data []byte) {
	Send(buf, 4, data)
}
