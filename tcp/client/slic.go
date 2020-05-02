// tut_tcpclient_filereceiver project slic.go
// Made by Gilles Van Vlasselaer
// More info about it on www.mrwaggel.be/post/golang-sending-a-file-over-tcp/

package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Print("usage: ./slic <dest> <tcp>")
		os.Exit(1)
	}
	connection, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer connection.Close()
	fc := make(chan string, 10)
	defer close(fc)
	if err := list(connection, fc, os.Args[1], "", os.Args[2]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go func() {
		for f := range fc {
			println(f)
		}
	}()
}

func list(connection net.Conn, fc chan string, src, prefix, dest string) error {
	fmt.Printf("listing... %v\n", src)
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		if err := filepath.Walk(src,  func(p string, info os.FileInfo, err error) error {
			if p == src {
				return nil
			}
			return list(connection, fc, p, filepath.Join(prefix, filepath.Dir(src)), dest)
		}); err != nil {
			return err
		}
	case mode.IsRegular():
		fc <- src
	}
	return nil
}