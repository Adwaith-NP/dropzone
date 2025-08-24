package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Adwaith-NP/dropzone/cmd"
	"github.com/Adwaith-NP/dropzone/internal/utils"
)

const DEFAULT_PORT int = 8080
const DEFAULT_NAME string = "DropZone"
const URL string = "/Users/adwaith/Documents/dropzone/cmd/receive.go"

func main() {
	// Get local ip and cheak is there any error
	localIp, err := utils.GetLocalIP()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error : ", err)
	}

	//get download directory

	downloadDir := utils.GetDownloadDir()
	if !utils.IsDirectory(downloadDir) {
		fmt.Fprintln(os.Stderr, "Error: Download directory not found. Please specify it manually using -dd <path>.")
		os.Exit(1)
	}

	// collect all command line argum
	senderMode := flag.Bool("s", false, "Run in sender mode")
	receiverMode := flag.Bool("r", false, "Run in receiver mode")

	ip := flag.String("ip", localIp, "Custom IPV4 setup")
	port := flag.Int("p", DEFAULT_PORT, "Custom port setup <warning : also set same port on reciever or sender>")
	dropName := flag.String("n", DEFAULT_NAME, "Custom drop name setup")
	filePath := flag.String("f", "", "File path")
	directoryPath := flag.String("d", "", "Directory path")
	fileList := flag.String("fl", "", "Comma-separated list of file paths")
	showHelp := flag.Bool("h", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")
	userDefinedDownloadDir := flag.String("dd", downloadDir, "Custom download directory")
	flag.Parse()

	if *receiverMode && *userDefinedDownloadDir != downloadDir && !utils.IsDirectory(*userDefinedDownloadDir) {
		fmt.Fprintln(os.Stderr, "Error: Directory not found : ", *userDefinedDownloadDir)
		os.Exit(1)
	}

	//validate user defined dropZone name
	if *dropName != DEFAULT_NAME {
		if len(*dropName) >= 15 {
			fmt.Fprintln(os.Stderr, "Error: Drop name can contain at most 15 letters")
			os.Exit(1)
		}
		if strings.Contains(*dropName, "|") {
			fmt.Fprintln(os.Stderr, "Error: Drop name cannot contain \"|\"")
			os.Exit(1)
		}
	}
	//Validate the given use given ip is valid
	if *ip != localIp && !utils.IsValidIPv4(*ip) {
		fmt.Fprintln(os.Stderr, "Error : Invalid IPV4: ", *ip)
		os.Exit(1)
	}
	//Display all command line argum
	if *showHelp {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	//Display version
	if *showVersion {
		fmt.Println("DropZone v1.0.0")
		os.Exit(0)
	}

	// check user gives both -s and -v
	if *senderMode && *receiverMode {
		fmt.Fprintln(os.Stderr, "Error : You cannot select both -s and -v")
		os.Exit(1)
	}

	// separate the login to sender mode or receiver mode by user specification
	if *senderMode {
		if *filePath != "" {
			cmd.SenderMode(*port, *filePath, utils.FileType, *ip)
		} else if *directoryPath != "" {
			cmd.SenderMode(*port, *directoryPath, utils.DirectoryType, *ip)
		} else if *fileList != "" {
			cmd.SenderMode(*port, *fileList, utils.FileListType, *ip)
		} else {
			fmt.Fprintln(os.Stderr, "Error : You must specify file(-f),directory(-d) or file list(-lf)")
			os.Exit(1)
		}
	} else if *receiverMode {
		cmd.ReceiverMode(*port, *dropName, *userDefinedDownloadDir, *ip)
	} else {
		fmt.Fprintln(os.Stderr, "Error: You must specify either -s (sender) or -r (receiver)")
		os.Exit(1)
	}

}
