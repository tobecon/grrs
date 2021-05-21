package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type FileInfo struct {
	FileName  string
	LineCount int32
}

var wg sync.WaitGroup

func main() {

	patthern := os.Args[1]
	dir := os.Args[2]
	entries, err := ioutil.ReadDir(dir)
	if err == nil {
		// for _, f := range entries {
		// 	if !f.IsDir() {
		// 		wg.Add(1)
		// 		go matchPattern(patthern, dir+"\\"+f.Name())
		// 	} else {

		// 	}
		// }
		err := workDir(entries, patthern, dir)
		if err != nil {
			fmt.Println(err)
		}
	}
	wg.Wait()
}

func workDir(entries []fs.FileInfo, patthern string, prefix string) (error error) {
	if len(entries) == 0 {
		return
	} else {
		for _, f := range entries {
			if !f.IsDir() {
				go matchPattern(patthern, prefix+"\\"+f.Name())
			} else {
				curPrefix := prefix + "\\" + f.Name()
				curentries, err := ioutil.ReadDir(prefix + "\\" + f.Name())
				if err == nil {
					workDir(curentries, patthern, curPrefix)
				} else {
					return err
				}
			}
		}
	}
	return
}

func matchPattern(pattern string, path string) (err error) {
	wg.Add(1)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	defer wg.Done()
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
