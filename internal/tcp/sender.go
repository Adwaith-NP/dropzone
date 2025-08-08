package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
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

	//Sending the len of json
	_, err = conn.Write(lengthBytes)
	if err != nil {
		return err
	}

	//sending the leg
	_, err = conn.Write(jsonData)
	if err != nil {
		return err
	}
	//read the response by the receiver
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	response := make([]byte, 8)
	n, err := conn.Read(response)
	if err != nil {
		return err
	}
	res := string(response[:n])

	if res == "accepted" {
		fmt.Println("Starting")
	} else {
		return fmt.Errorf("receiver rejected your request")
	}
	return nil
}
