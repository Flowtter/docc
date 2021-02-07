package main

import (
	"path"

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

func main() {
	var folder string = "html-docc"
	createFoldersIfDoesNotExist(folder)
	if !verifyFolders("include") {
		panic("no include folder")
	}

	var filesToCopy []string = []string{"prism.css", "prism.js", "style.css"}

	for i := 0; i < len(filesToCopy); i++ {
		copyFile(path.Join("assets", filesToCopy[i]), path.Join(folder, filesToCopy[i]))
	}

	var files []string = getAllFiles("include")
	var filesPath []string = make([]string, len(files))
	copy(filesPath, files)

	trimAllFiles("include/", files)
	trimAllFiles(".h", files)
	parseFiles(files, filesPath, folder)
	indexHTML(files, folder)
}
