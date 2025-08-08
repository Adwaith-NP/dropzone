package tcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/Adwaith-NP/dropzone/cmd"
)

func ReceiveMeta(post int) error {
	var meta map[string]any
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

	//Receive the leg of json
	lengthBytes := make([]byte, 8)
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return err
	}

	var length int64
	for i := 0; i < 8; i++ {
		length |= int64(lengthBytes[i]) << (8 * i)
	}

	//Receive json by using the len
	jsonBytes := make([]byte, length)
	if _, err := io.ReadFull(conn, jsonBytes); err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBytes, &meta); err != nil {
		return err
	}

	choice = cmd.RequestInquiry(meta)

	if choice == "y" {
		conn.Write([]byte("accepted"))
	} else {
		conn.Write([]byte("rejected"))
	}
	return nil
}
