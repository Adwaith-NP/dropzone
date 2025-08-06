package tcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

func ReceiveMeta(post int) error {
	var meta map[string]interface{}
	var choice string
	port := fmt.Sprintf(":%d", post)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	lengthBytes := make([]byte, 8)
	_, err = io.ReadFull(conn, lengthBytes)
	if err != nil {
		return err
	}
	var length int64
	for i := 0; i < 8; i++ {
		length |= int64(lengthBytes[i]) << (8 * i)
	}
	jsonBytes := make([]byte, length)
	_, err = io.ReadFull(conn, jsonBytes)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, &meta)
	if err != nil {
		return err
	}

	//Print the details of meta data
	fmt.Println("File transaction request (Accept or decline)")
	if meta["Type"] == "directory" {
		fmt.Printf("\nName : %s\nTotalSize : %d Byte\nFileCount : %d", meta["Name"], meta["TotalSize"], meta["FileCount"])
	} else {
		fmt.Printf("\nName : %s\nSize : %d Byte", meta["Name"], meta["Size"])
	}
	fmt.Printf("\n\n* Type n hit enter for reject\n* Type y and hit enter to accept\nInput : ")
	fmt.Scan(&choice)
	choice = strings.ToLower(choice)

	return nil
}
