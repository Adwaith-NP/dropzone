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
	fmt.Println("═════════════════════════════════════════════════════════════")
	fmt.Println("DropName           :", dropName)
	fmt.Printf("Your IP Address    : \033[32m%s\033[0m\n", ip)
	fmt.Printf("Listening on       : \033[36m%d\033[0m\n", port)
	fmt.Println("Download directory :", dowloadDir)
	fmt.Println("═════════════════════════════════════════════════════════════")

	//Using channel to stop and start broadcast
	startStopSignel := make(chan bool)

	//Start broadcasting by go goroutines
	go udp.StartBroadcast(dropName, port, ip, startStopSignel)
	//Start lisener for download acction
	err := tcp.ReceiveFiles(port, dowloadDir, startStopSignel)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error : ", err)
	}
}
