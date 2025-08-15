package cmd

import (
	"fmt"
	"os"

	"github.com/Adwaith-NP/dropzone/internal/tcp"
	"github.com/Adwaith-NP/dropzone/internal/udp"
	"github.com/Adwaith-NP/dropzone/internal/utils"
)

func ReceiverMode(port int, dropName string, dowloadDir string, ip string) {
	//Display details
	fmt.Print(utils.AsciiArt, "\n\n")
	fmt.Println("DropName           : ", dropName)
	fmt.Println("Your IPV4          : ", ip)
	fmt.Println("Used PORT          : ", port)
	fmt.Println("Download directory : ", dowloadDir)
	//Start broadcasting by go goroutines
	go udp.StartBroadcast(dropName, port)
	//Start lisener for download acction
	err := tcp.ReceiveFiles(port, dowloadDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error : ", err)
	}
}
