
/**
 * Auth :   liubo
 * Date :   2022/3/4 16:51
 * Comment: linux的tail工具
 */

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var f = flag.String("f", "", "-f=111.log")
var l = flag.Int("n", 10, "-n=10")


func main() {

	flag.Parse()

	if len(*f) == 0 {
		flag.PrintDefaults()
		panic("invalid filename")
		return
	}

	if *l == 0 {
		flag.PrintDefaults()
		panic("invalid args")
		return
	}

	var n = tailFile(*f, *l)
	fmt.Println("处理完毕, count=", n)
}

func removeFiles(patten string) {
	var dir = filepath.Dir(patten)
	var files, _ = ioutil.ReadDir(dir)
	for _, v := range files {

		if v.IsDir() {
			continue
		}

		if strings.HasPrefix(v.Name(), patten) {
			os.Remove(v.Name())
		}
	}
}

func tailFile(file string, lines int) int {
	var fileName = file
	var ext = filepath.Ext(file)
	if len(ext) == 0 {
		ext = ".log"
	}

	var tailFileName = func() string {
		return file + ".tail" + ext
	}

	// 删掉旧文件
	removeFiles(file + ".tail.")


	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return 0
	}

	defer fi.Close()

	var count = 0
	br := bufio.NewReader(fi)
	var bw = bytes.NewBuffer([]byte{})
	var lineData [][]byte = make([][]byte, 0, lines + 1)

	var addLineData = func(d []byte) {
		lineData = append(lineData, d)
		if len(lineData) == (lines + 1) {
			lineData = lineData[1:]
		}
	}

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			if len(a) > 0 {
				panic("extra line data")
				//addLineData(a)
			}
			break
		}
		addLineData(a)

		count++
	}

	for _, v := range lineData {
		bw.Write(v)
		bw.WriteString("\r\n")
	}

	ioutil.WriteFile(tailFileName(), bw.Bytes(), 666)

	return len(lineData)
}