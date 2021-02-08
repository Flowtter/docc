package function

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Flowtter/docc/utils"
	"github.com/gosimple/slug"
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

// Function in C language
type Function struct {
	Prototype   string
	Description string
	Path        string
	Line        int
}

func GetPostProcessedDescriptionOfFunction(endingLineIndex int, lines []string) string {
	return PostProcessingDescription(GetDescriptionOfFunction(endingLineIndex, lines))
}

func GetDescriptionOfFunction(endingLineIndex int, lines []string) string {
	index := endingLineIndex - 1
	if strings.Contains(lines[index], "//") {
		return lines[index]
	} else if strings.Contains(lines[index], "*/") {
		_, result := utils.ContainsStringInMultipleLines(index, "/*", lines)
		return result
	}
	return ""
}

func PostProcessingDescription(comment string) string {
	comment = strings.ReplaceAll(comment, "/**", "")
	comment = strings.ReplaceAll(comment, "/*", "")
	comment = strings.ReplaceAll(comment, "*/", "")
	split := strings.Split(comment, "\n")

	for i := 0; i < len(split); i++ {
		for len(split[i]) > 0 {
			if split[i][0] == ' ' || split[i][0] == '*' {
				split[i] = utils.TrimFirstRune(split[i])
			} else {
				break
			}
		}
	}

	// remove first lines
	for len(split) > 0 {
		if split[0] == "" {
			split = split[1:]
		} else {
			break
		}
	}

	//remove last lines
	for len(split) > 0 {
		if split[len(split)-1] == "" {
			split = split[:len(split)-1]
		} else {
			break
		}
	}

	return strings.Join(split, "\n")
}

func GetAllFunctionsOfLines(lines []string, path string) []Function {
	var functions []Function
	for index, line := range lines {
		var isComment = false
		if strings.Contains(line, "//") || strings.Contains(line, "/*") {
			for i := 0; i < len(line); i++ {
				if line[i] == ' ' {
					continue
				} else if line[i] == '/' {
					isComment = true
					break
				}
			}
		}

		if !isComment && strings.Contains(line, ");") {
			if strings.Contains(line, "(") {
				functions = append(functions, Function{
					Prototype:   line,
					Line:        index + 1,
					Path:        path,
					Description: GetPostProcessedDescriptionOfFunction(index, lines),
				})
			} else {
				startingLine, prototypeLong := utils.ContainsStringInMultipleLines(index, "(", lines)
				if startingLine == -1 {
					break
				}
				functions = append(functions, Function{
					Prototype:   prototypeLong,
					Line:        startingLine + 1,
					Path:        path,
					Description: GetPostProcessedDescriptionOfFunction(index, lines),
				})
			}
		}
	}
	return functions
}

func FolderMaker(path string) Folder {
	folder := Folder{
		Name: path,
	}
	folder.getFoldersRecursive(path)
	folder.First = true
	return folder
}

func (f *Folder) getFoldersRecursive(pth string) {
	all, err := ioutil.ReadDir(pth)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(all); i++ {
		if all[i].IsDir() {
			folder := Folder{
				Name:  all[i].Name(),
				First: false,
			}
			folder.getFoldersRecursive(path.Join(pth, all[i].Name()))
			f.SubFolders = append(f.SubFolders, folder)
		} else {
			joined := path.Join(pth, all[i].Name())
			fullName := utils.GetName([]string{joined})[0]
			f.Files = append(f.Files, File{
				Slug: slug.Make(all[i].Name())[:len(all[i].Name())-2],
				HREF: fullName + ".html",
				Path: joined,
			})
		}
	}
}

func ParseHTML(pageData PageData, pathToSave, pathFunctionFile string) {
	tmpl := template.Must(template.ParseFiles(path.Join("assets", "layout.html")))

	lines := utils.GetAllLinesOfFile(pathFunctionFile)
	if pageData.FolderTitle != "index" {
		pageData.Functions = GetAllFunctionsOfLines(lines, pathFunctionFile)
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

func ParseFiles(name, pathFunctionFile string, mainFolder Folder) {
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
	ParseHTML(pageData, path.Join(wd, "html-docc", name+".html"), pathFunctionFile)
}
