package main

import (
	"fmt"

	"github.com/Adwaith-NP/dropzone/internal/tcp"
	"github.com/Adwaith-NP/dropzone/internal/udp"
	"github.com/Adwaith-NP/dropzone/internal/utils"
)

const TCP_PORT = 8080
const UDP_PORT = 9090
const DEFAULT_NAME = "DropZone"
const URL = "/Users/adwaith/Documents/dropzone/cmd/receive.go"

func main() {
	test := false
	if test {
		ip, err := udp.StartListening(UDP_PORT)
		if err != nil {
			fmt.Println(err)
			return
		}
		meta, err := utils.BuildFileMeta(URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = tcp.SendFile(ip, TCP_PORT, meta, URL)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		go udp.StartBroadcast(DEFAULT_NAME, UDP_PORT)
		err := tcp.ReceiveFiles(TCP_PORT)
		if err != nil {
			fmt.Println("error ", err)
			return
		}
	}
}
