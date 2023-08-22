package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func IsExist(fn string) bool {
	_, err := os.Stat(fn)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}

	return false
}

func ReadDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	// 列出目录下所有文件名称
	for _, f := range files {
		fmt.Println(f.Name())
	}
	return nil
}

func ReadDirRecursive(dirPath string) error {
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// put your code here
			fmt.Println(path, info.Size())
			return nil
		})
	return err
}

func WriteFile() {
	path := "test.csv"

	// os.O_TRUNC 清空写
	// os.O_APPEND 追加写
	// os.O_CREATE 不存在则创建
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	write := bufio.NewWriter(f)
	for i := 0; i < 10; i++ {
		write.WriteString(fmt.Sprintf("%d,%s,%s\n",
			i,
			fmt.Sprintf("name-%d", i+1),
			time.Now().Add(-time.Duration(i)*time.Minute).Format("2006-01-02 15:04:05")),
		)
	}

	_ = write.Flush()
}

func SimpleWriteFile() {
	var output = ""
	for i := 0; i < 10; i++ {
		output += fmt.Sprintf("%d,%s,%s\n",
			i,
			fmt.Sprintf("name-%d", i+1),
			time.Now().Add(-time.Duration(i)*time.Minute).Format("2006-01-02 15:04:05"),
		)
	}
	err := os.WriteFile("test.csv", []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}

func SimpleReadFile() {
	data, err := os.ReadFile("test.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
