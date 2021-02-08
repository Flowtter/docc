package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)
}

// TODO: slugify path
// TODO: online (arg local)
// Todo: readme
func main() {

	args := os.Args[1:]

	var folder string = "html-docc"
	createFoldersIfDoesNotExist(folder)
	if !verifyFolders("include") {
		panic("no include folder")
	}

	var filesToCopy []string = []string{"prism.css", "prism.js", "style.css", "file-tree.css", "file-tree.js"}

	for i := 0; i < len(filesToCopy); i++ {
		copyFile(path.Join("assets", filesToCopy[i]), path.Join(folder, filesToCopy[i]))
	}

	var files []string = getAllFilesOrFolder("include", true)
	var filesName []string = getName(files)

	var mainFolder Folder = folderMaker("include")

	for i := 0; i < len(files); i++ {
		parseFiles(filesName[i], files[i], mainFolder)
	}

	parseFiles("index", "index.html", mainFolder)

	if len(args) > 0 {
		argsString := strings.Join(args, "")
		if strings.Contains(argsString, "help") {
			fmt.Println(`Welcome to docc,
	-help: displays help
	-l: launches the index.html`)
		} else if strings.Contains(argsString, "l") {
			openBrowser("html-docc/index.html")
		} else {
			fmt.Println("Unknown command maybe you should try using the help command")
		}
	}
}
