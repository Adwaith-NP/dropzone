package main

import (
	"fmt"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

const UDP_IP = 9090
const DEFAULT_NAME = "DropZone"

func main() {
	// udp.StartBroadcast(DEFAULT_NAME, UDP_IP)
	// ip, err := udp.StartListening(UDP_IP)
	// fmt.Println(ip, err)
	files, err := utils.GetAllFiles("/Users/adwaith/Documents/dropzone")
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range files {
		fmt.Println(i)
	}
}
