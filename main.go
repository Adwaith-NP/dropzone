package main

import (
	"fmt"

	"github.com/Adwaith-NP/dropzone/internal/udp"
)

const UDP_IP = 9090
const DEFAULT_NAME = "DropZone"

func main() {
	ip, err := udp.StartListening(UDP_IP)
	fmt.Println(ip, err)
}
