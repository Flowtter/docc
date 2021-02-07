package main

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/gosimple/slug"
)

// Function in C language
type Function struct {
	Prototype   string
	Description string
	Path        string
	Line        int
}

func getPostProcessedDescriptionOfFunction(endingLineIndex int, lines []string) string {
	return postProcessingDescription(getDescriptionOfFunction(endingLineIndex, lines))
}

func getDescriptionOfFunction(endingLineIndex int, lines []string) string {
	index := endingLineIndex - 1
	if strings.Contains(lines[index], "//") {
		return lines[index]
	} else if strings.Contains(lines[index], "*/") {
		_, result := containsStringInMultipleLines(index, "/*", lines)
		return result
	}
	return ""
}

func postProcessingDescription(comment string) string {
	comment = strings.ReplaceAll(comment, "/**", "")
	comment = strings.ReplaceAll(comment, "/*", "")
	comment = strings.ReplaceAll(comment, "*/", "")
	split := strings.Split(comment, "\n")

	for i := 0; i < len(split); i++ {
		for len(split[i]) > 0 {
			if split[i][0] == ' ' || split[i][0] == '*' {
				split[i] = trimFirstRune(split[i])
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

func getAllFunctionsOfLines(lines []string, path string) []Function {
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
					Description: getPostProcessedDescriptionOfFunction(index, lines),
				})
			} else {
				startingLine, prototypeLong := containsStringInMultipleLines(index, "(", lines)
				if startingLine == -1 {
					break
				}
				functions = append(functions, Function{
					Prototype:   prototypeLong,
					Line:        startingLine + 1,
					Path:        path,
					Description: getPostProcessedDescriptionOfFunction(index, lines),
				})
			}
		}
	}
	return functions
}

func folderMaker(path string) Folder {
	folder := Folder{
		Name: path,
	}
	folder.getFoldersRecursive(path)
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
				Name: all[i].Name(),
			}
			folder.getFoldersRecursive(path.Join(pth, all[i].Name()))
			f.SubFolders = append(f.SubFolders, folder)
		} else {
			joined := path.Join(pth, all[i].Name())
			f.Files = append(f.Files, File{
				Slug: slug.Make(all[i].Name())[:len(all[i].Name())-2],
				HREF: slug.Make(all[i].Name())[:len(all[i].Name())-2] + ".html",
				Path: joined,
			})
		}
	}
}
