package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	FileName  string
	LineCount int32
}

func main() {

	patthern := os.Args[1]
	dir := os.Args[2]
	entries, err := ioutil.ReadDir(dir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*400)
	defer cancel()
	errCh := make(chan int)
	go wordDirWithTimeOut(ctx, errCh, entries, patthern, dir)
	if err == nil {
		select {
		case <-ctx.Done():
			fmt.Println("end!")
		case i := <-errCh:
			fmt.Println(i)
		}
	}
}

func wordDirWithTimeOut(ctx context.Context, ch chan int, entries []fs.FileInfo, patthern string, prefix string) {
	workDir(ctx, ch, entries, patthern, prefix)
	ch <- 0
}

func workDir(ctx context.Context, ch chan int, entries []fs.FileInfo, patthern string, prefix string) {
	if len(entries) == 0 {
		//ch <- 0 // 为了持续阻塞 chan , 让程序不退出. 所以不需要在这个发信号  chan
	} else {
		for _, f := range entries {
			if !f.IsDir() {
				go matchPattern(patthern, prefix+"\\"+f.Name())
			} else {
				curPrefix := prefix + "\\" + f.Name()
				curentries, err := ioutil.ReadDir(prefix + "\\" + f.Name())
				if err == nil {
					workDir(ctx, ch, curentries, patthern, curPrefix)
				}
			}
		}
	}
}

func matchPattern(pattern string, path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)
	fileInfo := FileInfo{path, 0}
	for reader.Scan() {
		line := reader.Text()
		if strings.Contains(line, pattern) {
			fileInfo.LineCount++
		}
	}
	fileInfo.PrettyPrint()
	return nil
}
func (f *FileInfo) PrettyPrint() {
	if f.LineCount < 4 {
		return
	}
	fmt.Printf("file %v has %v line \n", f.FileName, f.LineCount)
}
