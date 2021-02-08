package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Flowtter/docc/function"
	"github.com/Flowtter/docc/utils"

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

// TODO: filetree reset
func main() {

	args := os.Args[1:]

	var folder string = "html-docc"
	utils.CreateFoldersIfDoesNotExist(folder)
	if !utils.VerifyFolders("include") {
		panic("no include folder")
	}

	var filesToCopy []string = []string{"prism.css", "prism.js", "style.css", "file-tree.css", "file-tree.js"}

	for i := 0; i < len(filesToCopy); i++ {
		utils.CopyFile(path.Join("assets", filesToCopy[i]), path.Join(folder, filesToCopy[i]))
	}

	var files []string = utils.GetAllFilesOrFolder("include", true)
	var filesName []string = utils.GetName(files)

	var mainFolder function.Folder = function.FolderMaker("include")

	for i := 0; i < len(files); i++ {
		function.ParseFiles(filesName[i], files[i], mainFolder)
	}

	function.ParseFiles("index", "index.html", mainFolder)

	if len(args) > 0 {
		argsString := strings.Join(args, "")
		if strings.Contains(argsString, "help") {
			fmt.Println(`Welcome to docc,
	-help: displays help
	-l: launches the index.html`)
		} else if strings.Contains(argsString, "l") {
			utils.OpenBrowser("html-docc/index.html")
		} else {
			fmt.Println("Unknown command maybe you should try using the help command")
		}
	}
}
