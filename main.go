package main

import (
	"encoding/json"
	"fmt"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

const UDP_PORT = 9090
const DEFAULT_NAME = "DropZone"

func main() {
	// udp.StartBroadcast(DEFAULT_NAME, UDP_IP)
	// ip, err := udp.StartListening(UDP_IP)
	// fmt.Println(ip, err)
	meta, err := utils.BuildDirectoryMeta("/Users/adwaith/Documents/dropzone")
	if err != nil {
		return
	}
	jsonBytes, _ := json.MarshalIndent(meta, "", "  ")
	fmt.Println(string(jsonBytes))
}
