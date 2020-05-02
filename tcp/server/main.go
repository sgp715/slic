// tut_tcpserver_filesend project slic.go
// Made by Gilles Van Vlasselaer
// More info about it on www.mrwaggel.be/post/golang-sending-a-file-over-tcp/

package main

import (
	"common"
	"fmt"
	"net"
	"os"
)

func main() {
	server, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("server started...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("client connected")
		go func() {
			if err := handleConnection(connection); err != nil {
				fmt.Printf("connection failed to copy %v", err)
			}
		}()
	}
}

func handleConnection(connection net.Conn) error {
	defer connection.Close()
	if err := common.ReceiveFile(connection); err != nil {
		return err
	}
	return nil
}