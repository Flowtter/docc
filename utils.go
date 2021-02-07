package main

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"unicode/utf8"
)

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
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

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}
}

func getName(files []string) []string {
	var name []string
	for i := 0; i < len(files); i++ {
		tmpName := files[i]
		_, tmpName = path.Split(tmpName)
		tmpName = tmpName[:len(tmpName)-2]
		name = append(name, tmpName)
	}
	return name
}
