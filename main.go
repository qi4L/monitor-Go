package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gookit/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var specialPathStr = "drops_B0503373BDA6E3C5CD4E5118C02ED13A"
var specialString = "drops_log"
var fileMD5Dict = make(map[string]string)
var originFileList []string

func calcMD5(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[*] 文件被删除 :", filePath)
		return ""
	}
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func getFilesList(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.Contains(path, specialPathStr) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("[*] 错误:", err)
	}
	return files
}

func backupDirectory(srcDir string) error {
	// 创建备份目录
	backupDir := filepath.Join(srcDir, "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// 读取源目录下的所有文件
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			srcFile := filepath.Join(srcDir, file.Name())
			dstFile := filepath.Join(backupDir, file.Name())

			if err := copyFile(srcFile, dstFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func replaceFileContent(targetFilePath, sourceFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, sourceFile)
	return err
}

func main() {
	fmt.Println("---------持续监测文件中------------")
	cwd, _ := os.Getwd()
	if err := backupDirectory(cwd); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Backup completed successfully!")
	}
	originFileList = getFilesList(cwd)
	for _, filePath := range originFileList {
		fileMD5Dict[filePath] = calcMD5(filePath)
	}

	for {
		fileList := getFilesList(cwd)
		diffFileList := []string{}
		for _, file := range fileList {
			if _, ok := fileMD5Dict[file]; !ok {
				diffFileList = append(diffFileList, file)
			}
		}

		// 移除新上传文件
		if len(diffFileList) != 0 {
			for _, filePath := range diffFileList {
				data, err := os.ReadFile(filePath)
				if err != nil {
					break
				}
				content := string(data)
				if !strings.Contains(content, specialString) {
					fmt.Println("[*] 发现疑似WebShell上传文件:", filePath, "时间为:", time.Now(), "内容为:", content)
					// 自动删除新上传的文件
					err1 := os.Remove(filePath)
					if err1 != nil {
						color.Red.Printf("删除文件时出错：%s\n", err)
					}
					color.Green.Println("[+] 新上传的文件已删除")
				}
			}
		}

		// 防止任意文件被修改,还原被修改文件
		for _, filePath := range originFileList {
			newMD5 := calcMD5(filePath)
			if newMD5 != fileMD5Dict[filePath] {
				data, err1 := os.ReadFile(filePath)
				if err1 != nil {
					break
				}
				content := string(data)
				if !strings.Contains(content, specialString) {
					fmt.Println("[*] 该文件被修改 :", filePath, "时间为:", time.Now(), "内容为:", content)
					// 自动还原被删除的文件
					dir, file := filepath.Split(filePath)
					backupFilePath := filepath.Join(dir, "backup", file)
					replaceFileContent(filePath, backupFilePath)
					color.Green.Println("[+] 被修改的文件已还原")
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
