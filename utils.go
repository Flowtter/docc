package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Function struct {
	Prototype   string
	Description string
	Path        string
	Line        int
}

func printFunction(fun Function) {
	fmt.Println("Function in", fun.Path)
	fmt.Println("At line", fun.Line)
	fmt.Println(fun.Prototype)
	fmt.Println(fun.Description)
	fmt.Println()
}

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

func getAllLinesOfFile(path string) []string {
	file, _ := ioutil.ReadFile(path)
	return strings.Split(strings.ReplaceAll(string(file), "\r", ""), "\n")
}

func getAllFiles(folder string) []string {
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !isFolder(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func isFolder(folder string) bool {
	file, err := os.Stat(folder)
	if err != nil {
		return false
	}
	return file.Mode().IsDir()
}

// debug
func getAllFilesAndTrim(folder string) []string {
	files := getAllFiles(folder)
	trimAllFiles(folder, files)
	return files
}

func trimAllFiles(trim string, files []string) {
	for i := 0; i < len(files); i++ {
		files[i] = strings.ReplaceAll(files[i], trim, "")
	}
}

func containsStringInMultipleLines(endingLineIndex int, startingString string, lines []string) (int, string) {
	resultString := lines[endingLineIndex]

	for index := endingLineIndex - 1; index > 0; index-- {
		resultString = lines[index] + "\n" + resultString
		if strings.Contains(lines[index], startingString) {
			return index, resultString
		}
	}
	return -1, ""
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

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// delete empty elements
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

// NETWORK : network_sgd
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
