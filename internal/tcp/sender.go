package tcp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// Start at the metaData sending to user for better communication
func SendFile(ip string, port int, meta any, url string) error {
	sendTo := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", sendTo)
	if err != nil {
		return err
	}
	defer conn.Close()
	jsonData, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	length := int64(len(jsonData))
	lengthBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		lengthBytes[i] = byte(length >> (8 * i))
	}

	//Sending the len of json
	_, err = conn.Write(lengthBytes)
	if err != nil {
		return err
	}

	//sending the leg
	_, err = conn.Write(jsonData)
	if err != nil {
		return err
	}
	//read the response by the receiver
	fmt.Println("waiting for receiver to accept ...")
	response := make([]byte, 8)
	n, err := conn.Read(response)
	if err != nil {
		return err
	}
	res := string(response[:n])

	if res != "accepted" {
		return fmt.Errorf("receiver rejected your request")
	}
	err = sendFiles(conn, url)
	return nil
}

// it a two way that , as per file or directory the path was selected
func sendFiles(conn net.Conn, url string) error {
	if utils.IsDirectory(url) {
		sendDirectory(conn, url)
	} else if utils.PathExists(url) {
		sendSingleFiles(conn, url, filepath.Base(url))
	}

	return nil
}

// add the logic to get the file url and pass to sendSingleFiles
func sendDirectory(conn net.Conn, url string) error {
	if !utils.IsDirectory(url) {
		return fmt.Errorf("invalid directory URL")
	}

	allFilesUrl, err := utils.GetAllFiles(url)
	if err != nil {
		return err
	}
	for _, file := range allFilesUrl {
		rel, err := filepath.Rel(filepath.Dir(url), file)
		if err != nil {
			return err
		}
		fmt.Println(url, file, rel)
		if err := sendSingleFiles(conn, file, rel); err != nil {
			fmt.Println("file problem : ", err)
		}
	}
	return nil
}

// send one file at a time , look up for all possible error
func sendSingleFiles(conn net.Conn, file string, url_name string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	meta := map[string]string{
		"path": url_name,
		"size": strconv.FormatInt(info.Size(), 10),
	}

	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	length := int64(len(metaBytes))
	if err := binary.Write(conn, binary.LittleEndian, length); err != nil {
		return err
	}

	if _, err := conn.Write(metaBytes); err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(conn, f); err != nil {
		return err
	}
	return nil
}
