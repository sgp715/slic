package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"flag"
	"golang.org/x/sync/errgroup"
	"strings"
)

var src = flag.String("src", "", "directory to upload")
var dest = flag.String("dest", "", "directory to send files")
var host = flag.String("host", "", "host to send files")
var jobs = flag.Int("jobs", 100, "max number of jobs to use")


func main() {
	flag.Parse()
	fi, err := os.Stat(*src)
	if *jobs < 1 {
		fmt.Println("jobs >= 1")
		os.Exit(1)
	}
	if err != nil {	
		fmt.Println("src invalid")
		os.Exit(1)
	}
	if !strings.Contains(*host, "@") {
		fmt.Println("host invalid")
		os.Exit(1)
	}
	if !fi.Mode().IsDir() {
		fmt.Println("src must be dir")
		os.Exit(1)
	}
	if err := syncDirs(*src, *dest, *host); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fc := make(chan fileDest, *jobs)
	var g errgroup.Group
	g.Go(func() error {
		err := list(fc, *src)
		close(fc)
		return err
	})
	for f := range fc {
		fmt.Printf("copying: %v\n", f)
		aPath := f.aPath
		g.Go(func() error {
			return scp(aPath, f.dest(*host, *dest))
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func syncDirs(src, dest, host string) error {
	fmt.Printf("sending %v -> %v\n", src, dest)
	if err := filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if err := rmkdir(host, filepath.Join(dest, strings.TrimPrefix(p, src))); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

type fileDest struct {
	file string
	aPath string
	rPath string
}

func (fd fileDest) dest(host, dest string) string {
	return fmt.Sprintf("%v:%v", host, filepath.Join(dest, fd.rPath))
}

func list(fc chan fileDest, src string) error {
	if err := filepath.Walk(src, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fc <- fileDest{ file: filepath.Base(p), aPath: p, rPath: strings.TrimPrefix(p, src) }
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func scp(path, dest string) error {
	scp := fmt.Sprintf("scp %v %v", path, dest)
	fmt.Printf("executing: %v\n", scp)
	cmd := exec.Command("bash", "-c", scp)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

func rmkdir(host, dest string) error {
	mkdir := fmt.Sprintf("ssh %v \"mkdir %v\"", host, dest)
	fmt.Printf("executing: %v\n", mkdir)
	cmd := exec.Command("bash", "-c", mkdir)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil && !strings.Contains(string(out),"exists") {
		return err
	}
	return nil
}