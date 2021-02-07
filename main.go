package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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

// TODO: folders
// TODO: online (arg local)
// Todo: readme
func main() {

	err := filepath.Walk("include",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	panic("hihi")

	args := os.Args[1:]

	var folder string = "html-docc"
	createFoldersIfDoesNotExist(folder)
	if !verifyFolders("include") {
		panic("no include folder")
	}

	var filesToCopy []string = []string{"prism.css", "prism.js", "style.css"}

	for i := 0; i < len(filesToCopy); i++ {
		copyFile(path.Join("assets", filesToCopy[i]), path.Join(folder, filesToCopy[i]))
	}

	var files []string = getAllFilesOrFolder("include", true)
	var filesPath []string = make([]string, len(files))
	copy(filesPath, files)

	trimAllFiles("include/", files)
	trimAllFiles(".h", files)
	parseFiles(files, filesPath, folder)

	fmt.Println(getAllFilesOrFolder("include", false))

	indexHTML(files, folder)

	if len(args) > 0 {
		argsString := strings.Join(args, "")
		if strings.Contains(argsString, "help") {
			fmt.Println(`Welcome to docc,
	-help: displays help
	-l: launches the index.html
	-e: opens the index.html`)
		} else if strings.Contains(argsString, "l") {
			openBrowser("html-docc/index.html")
			// } else if strings.Contains(argsString, "e") {

			// 	wd, err := os.Getwd()

			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	fmt.Println("explorer.exe", path.Join(wd, "html-doc"))
			// 	cmd := exec.Command("explorer.exe " + path.Join(wd, "html-doc"))

			// 	err = cmd.Run()

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
		} else {
			fmt.Println("Unknown command maybe you should try using the help command")
		}
	}
}
