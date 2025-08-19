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
	"strings"
	"time"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// Start at the metaData sending to user for better communication
func RequestInquiry(ip string, port int, meta any) (net.Conn, error) {
	sendTo := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", sendTo)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(meta) //Convert meta data to json
	if err != nil {
		conn.Close()
		return nil, err
	}
	length := uint64(len(jsonData))                                         //store length of json to uint64
	if err := binary.Write(conn, binary.LittleEndian, length); err != nil { //convert uint64 length to byte array and send to receiver
		conn.Close()
		return nil, err
	}

	//sending the leg
	_, err = conn.Write(jsonData) //sending jsonData
	if err != nil {
		conn.Close()
		return nil, err
	}
	//read the response by the receiver
	fmt.Print("\n\nWaiting for receiver to accept ...")
	response := make([]byte, 8)
	n, err := conn.Read(response)
	if err != nil {
		conn.Close()
		return nil, err
	}
	res := string(response[:n])

	if res == utils.StatusRejected {
		conn.Close()
		return nil, fmt.Errorf("receiver rejected your request")
	}
	return conn, nil
}

// add the logic to get the file url and pass to sendSingleFiles
func SendDirectory(conn net.Conn, url string) error {

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
		if err := SendSingleFiles(conn, file, rel); err != nil {
			fmt.Println("Error in upload: ", err)
		}
	}
	return nil
}

func SendFileList(conn net.Conn, fileList []string) {
	for _, file := range fileList {
		fileName := filepath.Base(file)
		if err := SendSingleFiles(conn, file, fileName); err != nil {
			fmt.Println("Error in upload: ", err)
		}
	}
}

// send one file at a time , look up for all possible error
func SendSingleFiles(conn net.Conn, file string, urlForDir string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	meta := map[string]string{
		"path": urlForDir,                          //Its use to determine where the file want to add when sharing mode on directory.
		"size": strconv.FormatInt(info.Size(), 10), //Size of the file in byte
	}

	metaBytes, err := json.Marshal(meta) //Converting meta data to json
	if err != nil {
		return err
	}

	length := uint64(len(metaBytes))                                        //Get the length of meta data
	if err := binary.Write(conn, binary.LittleEndian, length); err != nil { //convert the uint64 length to byte array
		return err
	}

	if _, err := conn.Write(metaBytes); err != nil { //Send the byte array
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	wd := &utils.DropData{
		Writer:        conn,
		LastTime:      time.Now(),
		TotalFileSize: info.Size(),
	}

	fmt.Printf("\n\nSending file: %s\n", info.Name())

	if _, err := io.Copy(wd, f); err != nil { //send file
		return err
	}
	fmt.Printf("\r[%s] 100%% (Done)\n", strings.Repeat("|", utils.BARWIDTH))
	return nil
}
