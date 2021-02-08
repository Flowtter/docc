package main

import (
	"os"
	"path"
	"strings"
	"text/template"
)

// PageData struct for html info
type PageData struct {
	PageTitle   string
	FolderTitle string
	MainFolder  Folder
	Functions   []Function
}

// Folder for the left menu
type Folder struct {
	Name       string
	Files      []File
	SubFolders []Folder
	First      bool
}

// File for the left menu
type File struct {
	Path string
	HREF string
	Slug string
}

func parseHTML(pageData PageData, pathToSave, pathFunctionFile string) {
	tmpl := template.Must(template.ParseFiles(path.Join("assets", "layout.html")))

	lines := getAllLinesOfFile(pathFunctionFile)
	if pageData.FolderTitle != "index" {
		pageData.Functions = getAllFunctionsOfLines(lines, pathFunctionFile)
	}
	f, err := os.Create(pathToSave)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, pageData)
	if err != nil {
		panic(err)
	}
}

func parseFiles(name, pathFunctionFile string, mainFolder Folder) {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	_, pageTitle := path.Split(wd)

	pageData := PageData{
		PageTitle:   strings.ToUpper(pageTitle),
		FolderTitle: pathFunctionFile,
		MainFolder:  mainFolder,
	}
	parseHTML(pageData, path.Join(wd, "html-docc", name+".html"), pathFunctionFile)
}
