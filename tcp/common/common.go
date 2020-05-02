package common

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 65495

func fillString(retunString string, toLength int) string {
	for {
		lengthString := len(retunString)
		if lengthString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func SendFile(connection net.Conn, src, dest string) error {
	fmt.Printf("sending file %v\n", src)
	file, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(dest, 64)
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Printf("sent file %v\n", src)
	return nil
}
func ReceiveFile(connection net.Conn) error {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	connection.Read(bufferFileSize)
	fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	if err != nil {
		return err
	}
	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")
	newFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("could not create file %v", err)
	}
	defer newFile.Close()
	fmt.Printf("writing file %v\n", fileName)
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Printf("wrote file %v\n", fileName)
	return nil
}
