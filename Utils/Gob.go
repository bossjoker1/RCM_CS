package Utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func Serialize(data []string) []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	// 将b编码后存进res输出流
	if err := encoder.Encode(data); err != nil {
		log.Printf("serialize the block to []byte failed %v\n", err)
	}
	return res.Bytes()
}

func Deserialize(databytes []byte) []string {
	var data []string
	decoder := gob.NewDecoder(bytes.NewReader(databytes))
	// 从输入流中读取进b
	if err := decoder.Decode(&data); err != nil {
		log.Printf("deserialize the []byte to []string failed. %v\n", err)
	}
	return data
}
