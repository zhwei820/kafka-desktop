package conf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func saveFileWithHistory(sourceFilePath, destinationDir string, maxHistoryOpt ...int) error {
	maxHistory := 5
	if len(maxHistoryOpt) > 0 {
		maxHistory = maxHistoryOpt[0]
	}
	// 获取文件名和扩展名
	sourceFileName := filepath.Base(sourceFilePath)
	fileExt := filepath.Ext(sourceFileName)
	fileNameWithoutExt := strings.TrimSuffix(sourceFileName, fileExt)

	// 获取当前时间作为时间戳
	timestamp := time.Now().Unix()

	// 构建目标文件名
	destinationFileName := fileNameWithoutExt + "_" + strconv.FormatInt(timestamp, 10) + fileExt
	destinationPath := filepath.Join(destinationDir, destinationFileName)

	// 检查目标目录是否存在，如果不存在则创建
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destinationDir, 0755); err != nil {
			return err
		}
	}

	// 复制源文件到目标位置
	if err := copyFile(sourceFilePath, destinationPath); err != nil {
		return err
	}

	for {
		// 获取目标目录下的所有历史文件
		files, err := filepath.Glob(filepath.Join(destinationDir, fileNameWithoutExt+"_*"+fileExt))
		if err != nil {
			return err
		}
		fmt.Println("len(files)", len(files))
		// 如果历史文件数量超过最大历史文件数，删除最旧的文件
		if len(files) > maxHistory {
			oldestFile := findOldestFile(files)
			if oldestFile != "" {
				if err := os.Remove(oldestFile); err != nil {
					return err
				}
			}
		} else {
			break
		}
	}

	return nil
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func findOldestFile(files []string) string {
	var oldestFile string
	var oldestTime time.Time

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		if oldestFile == "" || fileInfo.ModTime().Before(oldestTime) {
			oldestFile = file
			oldestTime = fileInfo.ModTime()
		}
	}

	return oldestFile
}
