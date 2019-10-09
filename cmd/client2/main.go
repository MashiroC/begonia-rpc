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
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/sdk"
	"time"
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

	pubReq := begonia.Request{
		Service:  "key",
		Function: "Public",
		//Param:    entity.Param{"info": entity.Param{"name":"hc"}, "sub": "test","duration":time.Hour * 3},
	}

	pubPki := cli.Call(pubReq).String()
	fmt.Println(pubPki)

	req := begonia.Request{
		Service:  "key",
		Function: "CreateToken",
		Param:    entity.Param{"info": entity.Param{"name": "hc"}, "sub": "test", "duration": time.Hour * 3},
	}
	resp := cli.Call(req)
	fmt.Println(resp.Error())
	fmt.Println(resp.Data)
	token := resp.String()

	//arr := strings.Split(token, ".")
	fmt.Println(token)
	fmt.Println(resp.Error())
	//fmt.Println()
	//fmt.Println(CheckToken(arr[0], arr[1],pubPki))
	//fmt.Println(resp.Error())
	//if resp.Error() != nil {
	//	fmt.Println(resp.Error())
	//} else {
	//	fmt.Println(i)
	//}
	//
	//cli.CallAsync(req, func(resp *begonia.Response) {
	//	i := resp.Int()
	//	if resp.Error() != nil {
	//		fmt.Println(resp.Error())
	//	} else {
	//		fmt.Println(i)
	//	}
	//})

	//wait := make(chan bool)
	//<-wait
}
