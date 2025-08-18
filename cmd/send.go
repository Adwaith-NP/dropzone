package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Adwaith-NP/dropzone/internal/tcp"
	"github.com/Adwaith-NP/dropzone/internal/udp"
	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// The step by step sender procces
func SenderMode(port int, path string, pathMode string, localIp string) {
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

	switch pathMode {
	case utils.FileType: //Logic then path is a file
		//Get full path
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting absolute path: ", err)
			os.Exit(1)
		}
		if utils.IsDirectory(absPath) { //Check the given path exist
			fmt.Fprintln(os.Stderr, "File path not found")
			os.Exit(1)
		}
		fileMeta, err := utils.BuildFileMeta(absPath) //Get the meta data of the file
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to generate meta data: ", err)
			os.Exit(1)
		}
		conn, err := tcp.RequestInquiry(receiverIp, port, fileMeta) //Send request to receiver if reject the exicution stop
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when connecting: ", err)
			os.Exit(1)
		}
		defer conn.Close()
		if err = tcp.SendSingleFiles(conn, path, filepath.Base(path)); err != nil { //filepath.Base(path) give the file name from the given url
			fmt.Fprintln(os.Stderr, "Error when sending file: ", err)
			os.Exit(1)
		} else {
			fmt.Print("\n\nFile sent successfully")
		}

	case utils.DirectoryType: //Logic then path is a directory
		//Get full path
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting absolute path: ", err)
			os.Exit(1)
		}
		if !utils.IsDirectory(absPath) { //Chack is dir exist
			fmt.Fprintln(os.Stderr, "Directory path not found")
			os.Exit(1)
		}
		dirMeta, err := utils.BuildDirectoryMeta(absPath) //create meta data of dir
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error to generate meta data: ", err)
			os.Exit(1)
		}
		conn, err := tcp.RequestInquiry(receiverIp, port, dirMeta) //Send request to receiver if reject the exicution stop
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when connecting: ", err)
			os.Exit(1)
		}
		defer conn.Close()
		if err = tcp.SendDirectory(conn, absPath); err != nil { //send directory
			fmt.Fprintln(os.Stderr, "Error when sending file: ", err)
			os.Exit(1)
		} else {
			fmt.Println("Directory sended")
		}
	case utils.FileListType:
		var totalSize int64 = 0
		fileNameSize := make(map[string]any)
		files := strings.Split(path, ",")
		if len(files) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No files specified. Please provide at least one file using -fl <file1,file2>")
			os.Exit(1)
		}
		for i := range files {
			file := strings.TrimSpace(files[i])
			absPath, err := filepath.Abs(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting absolute path: ", err)
				os.Exit(1)
			}
			if utils.IsDirectory(absPath) { //Chack is dir exist
				fmt.Fprintln(os.Stderr, "File path not found: ", files[i])
				os.Exit(1)
			}
			info, err := os.Stat(absPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "A problem in file: ", files[i], " Error: ", err)
				os.Exit(1)
			}
			totalSize += info.Size()
			fileNameSize[info.Name()] = info.Size()
			files[i] = absPath
		}

		meta := utils.BuildFileListMeta(fileNameSize, totalSize)
		conn, err := tcp.RequestInquiry(receiverIp, port, meta) //Send request to receiver if reject the exicution stop
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error when connecting: ", err)
			os.Exit(1)
		}
		defer conn.Close()
		tcp.SendFileList(conn, files)

	}
}
