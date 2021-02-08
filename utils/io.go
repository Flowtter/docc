package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CreateFoldersIfDoesNotExist(folder string) bool {
	if VerifyFolders(folder) {
		return true
	}
	return createFolders(folder)
}

func VerifyFolders(folder string) bool {
	_, err := os.Stat(folder)
	return !os.IsNotExist(err)
}

func createFolders(folder string) bool {
	return os.Mkdir(folder, 0755) == nil
}

func CopyFile(src, dst string) {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		panic(err)
	}
}

func IsFolder(folder string) bool {
	file, err := os.Stat(folder)
	if err != nil {
		return false
	}
	return file.Mode().IsDir()
}

func GetAllLinesOfFile(path string) []string {
	file, _ := ioutil.ReadFile(path)
	return strings.Split(strings.ReplaceAll(string(file), "\r", ""), "\n")
}

func GetAllFilesOrFolder(folder string, wantFile bool) []string {
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if wantFile == !IsFolder(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
