package utils

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
)

func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func ContainsStringInMultipleLines(endingLineIndex int, startingString string, lines []string) (int, string) {
	resultString := lines[endingLineIndex]

	for index := endingLineIndex - 1; index > 0; index-- {
		resultString = lines[index] + "\n" + resultString
		if strings.Contains(lines[index], startingString) {
			return index, resultString
		}
	}
	return -1, ""
}

func OpenBrowser(url string) {
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

func GetName(files []string) []string {
	var name []string
	for i := 0; i < len(files); i++ {
		tmpName := files[i]
		tmpFullName := strings.Split(tmpName, "/")
		tmpFullName = tmpFullName[1:]
		tmpName = path.Join(tmpFullName...)
		tmpName = tmpName[:len(tmpName)-2]
		name = append(name, slug.Make(tmpName))
	}
	return name
}
