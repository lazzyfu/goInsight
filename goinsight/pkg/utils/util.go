package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/alexmullins/zip"
)

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func EncryptAndTarGzFiles(inputFile, outputFile, filePath, key string) error {
	// 创建zip文件
	fzip, err := os.Create(fmt.Sprintf("%s/%s", filePath, outputFile))
	if err != nil {
		return err
	}
	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	// 打开待压缩的文件
	f, err := os.Open(fmt.Sprintf("%s/%s", filePath, inputFile))
	if err != nil {
		return err
	}
	defer f.Close()

	// 加密文件
	w, err := zipw.Encrypt(inputFile, key)
	if err != nil {
		return err
	}

	// 将文件内容复制到zip文件中
	if _, err = io.Copy(w, f); err != nil {
		return err
	}
	zipw.Flush()
	return nil
}
