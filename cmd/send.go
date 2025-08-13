package cmd

import (
	"fmt"
	"os"

	"github.com/Adwaith-NP/dropzone/internal/udp"
	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// The step by step sender procces
func SenderMode(port int, path string, patMode string, localIp string) {
	// Display ascii art of dropzone , local ip , and given path
	fmt.Print(utils.AsciiArt, "\n\n")
	fmt.Println("Your IP : ", localIp)
	fmt.Println("Path    : ", path)

	// Start searching for receivers
	receiverIp, err := udp.StartListening(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	fmt.Println(receiverIp)

	if patMode == utils.FileType {
	}
}
