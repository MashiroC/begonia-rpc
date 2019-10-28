// Time : 2019/8/22 下午9:30
// Author : MashiroC

// client2 something
package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	begonia "github.com/MashiroC/begonia-rpc/sdk"
)

func CheckToken(payload, signature, pubPki string) bool {

	// create PublicKey
	block, _ := pem.Decode([]byte(pubPki))
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	// base64 decode
	pl, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return false
	}

	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// hash payload
	h := sha256.New()
	h.Write(pl)
	hash := h.Sum(nil)
	// verify
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash, sign)
	return err == nil

}

func main() {
	cli := begonia.New(":4949")

	helloService := cli.Service("Hello")

	hello := helloService.Fun("Hello")

	res, err := hello("mashiroc")

	fmt.Println(res,err)
}

//func main() {
//	conn, err := net.Dial("tcp4", ":4949")
//	if err!=nil{
//		log.Fatal(err.Error())
//	}
//	c:=conn2.New(conn)
//	r:=entity.Request{
//		UUID:    "11111",
//		Service: "Hello",
//		Fun:     "Hello",
//		Data:    []interface{}{"hc"},
//	}
//	b, err := json.Marshal(r)
//	if err!=nil{
//		log.Fatal(err.Error())
//	}
//	_ = c.WriteRequest(b)
//
//	opcode,data,err:=c.ReadData()
//	if err!=nil{
//		log.Fatal(err.Error())
//	}
//	fmt.Println(opcode)
//	fmt.Println(string(data))
//	wait:=make(chan bool)
//	<-wait
//}
