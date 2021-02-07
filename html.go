package main

import (
	"html/template"
	"os"
	"path"
	"strings"
)

type PageData struct {
	PageTitle   string
	FolderTitle string
	PathFile    string
	Folders     []Folder
	Functions   []Function
}

type Folder struct {
	Name     string
	NameHTML string
}

func parseHTML(pageData PageData, wd string) {
	tmpl := template.Must(template.ParseFiles(path.Join("assets", "layout.html")))
	// todo : explode path
	pathFile := strings.ReplaceAll(pageData.FolderTitle+".html", "/", "-")

	lines := getAllLinesOfFile(pageData.PathFile)
	if pageData.FolderTitle != "index" {
		pageData.Functions = getAllFunctionsOfLines(lines, pageData.PathFile)
	}
	f, err := os.Create(path.Join(wd, pathFile))
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, pageData)
	if err != nil {
		panic(err)
	}
}

func parseFiles(files, filesPath []string, folderPath string) {
	var folders []Folder
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	_, pageTitle := path.Split(wd)

	folders = getFolders(files)

	wd = path.Join(wd, folderPath)

	for i := 0; i < len(files); i++ {
		pageData := PageData{
			PageTitle:   strings.ToUpper(pageTitle),
			FolderTitle: folders[i].Name,
			Folders:     folders,
			PathFile:    filesPath[i],
		}
		parseHTML(pageData, wd)
	}
}

func indexHTML(files []string, folderPath string) {
	var folders []Folder
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	_, pageTitle := path.Split(wd)

	for i := 0; i < len(files); i++ {
		folders = append(folders, Folder{
			Name:     files[i],
			NameHTML: strings.ReplaceAll(files[i]+".html", "/", "-"),
		})
	}

	wd = path.Join(wd, folderPath)

	pageData := PageData{
		PageTitle:   strings.ToUpper(pageTitle),
		FolderTitle: "index",
		Folders:     folders,
		PathFile:    "index.html",
	}
	parseHTML(pageData, wd)
}
