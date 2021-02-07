package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func createFoldersIfDoesNotExist(folder string) bool {
	if verifyFolders(folder) {
		return true
	}
	return createFolders(folder)
}

func verifyFolders(folder string) bool {
	_, err := os.Stat(folder)
	return !os.IsNotExist(err)
}

func createFolders(folder string) bool {
	return os.Mkdir(folder, 0755) == nil
}

func copyFile(src, dst string) {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		panic(err)
	}
}

func isFolder(folder string) bool {
	file, err := os.Stat(folder)
	if err != nil {
		return false
	}
	return file.Mode().IsDir()
}

func getAllLinesOfFile(path string) []string {
	file, _ := ioutil.ReadFile(path)
	return strings.Split(strings.ReplaceAll(string(file), "\r", ""), "\n")
}

func getAllFilesOrFolder(folder string, wantFile bool) []string {
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if wantFile == !isFolder(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
