package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var successCount, existsCount, errorCount = 0, 0, 0

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func CheckDir(path string) {
	if !Exist(path) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("文件夹夹创建失败:", err)
			log.Fatal(err)
		}
	}
}

func getFileName(url string, dirName string) string {
	urlList := strings.Split(url, "/")
	return dirName + "/" + urlList[len(urlList)-1]
}

func SaveFile(url string, dirName string, header map[string]string, delay int, wg *sync.WaitGroup) {
	defer wg.Done()
	filename := getFileName(url, dirName)
	if Exist(filename) {
		existsCount++
		printLog()
		return
	}
	resp, err := Get(url, header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	pix, _ := ioutil.ReadAll(resp.Body)
	if err := ioutil.WriteFile(filename, pix, 0777); err != nil {
		fmt.Println(err)
		errorCount++
		printLog()
		return
	}
	successCount++
	if delay > 0 {
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
	printLog()
}

func printLog() {
	fmt.Printf("Successed:%v, Existed:%v, Failure:%v\r", successCount, existsCount, errorCount)
}

func RandInt(min int, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max-min) + min)
}