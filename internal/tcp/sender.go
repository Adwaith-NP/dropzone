package tcp

import (
	"encoding/json"
	"fmt"
	"net"
)

func SendMetaData(ip string, port int, meta any) error {
	sendTo := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", sendTo)
	if err != nil {
		return err
	}
	defer conn.Close()
	jsonData, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	length := int64(len(jsonData))
	lengthBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		lengthBytes[i] = byte(length >> (8 * i))
	}
	_, err = conn.Write(lengthBytes)
	return err
}
