package main

import "github.com/Adwaith-NP/dropzone/internal/udp"

func main() {
	udp.StartBroadcast("adwaith", 9090)
}
