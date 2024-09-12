package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getClosingObjTag(strUICObj string) string {
	var tempArr []string

	tempArr = strings.SplitAfterN(strUICObj, ` `, 1)
	strRes := strings.Replace(tempArr[0], `<`, ``)

	return (`</` + strRes + `>`)
}

/*
------------------------------------------------------------------------------------------

funcname: SeachInFile Function
return: strings	: if found return "Line: %d" with %d as location line number

			  if not found return empty ""
	err		: return opening function if fail to open the file

args:	filePath	: file path

	searchProp	: property tag to search
	searchObj	: object tag to search

-------------------------------------------------------------------------------------------
*/
func searchInFile(filePath, searchProp, searchObj string) (string, error) {
	line := 1
	bCandFound := false

	var strBuilderObj strings.Builder

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strTempObj := scanner.Text()
		switch bCandFound {

		// If candidate is found
		case true:
			{
				if (strings.Contains(strTempObj, searchProp)) || (strings.Contains(strTempObj, getClosingObjTag(searchObj))) {
					bCandFound = false
					// If property was found inside the object tag, or the object tag ends
					// If property was found inside, assign result
					if strings.Contains(strTempObj, searchObj) {
						strResTmp := "Line: " + strconv.Itoa(line) + ","
						strBuilderObj.WriteString(strResTmp)
					}
				}
				break
			}

		// Default search for object
		default:
			{
				if strings.Contains(strTempObj, searchObj) {
					// Check if the object has closing tags, therefore properties are not overriden
					if strings.HasPrefix(strTempObj, "<") && strings.HasSuffix(strTempObj, "/>") {
						// do nothing
					} else {
						bCandFound = true
					}
				}
				break
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
	}

	return strBuilderObj.String(), nil
}

func searchInDirectory(dirPath, searchProp, searchObj string) error {

	err := filepath.WalkDir(dirPath, func(filePath string, info os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error walking the path %s: %v\n", filePath, err)
			return err
		}
		if !info.IsDir() {
			strSearchRes, err := searchInFile(filePath, searchProp, searchObj)
			if strSearchRes != "" || err != nil {
				fmt.Printf("String '%s' found in file: %s %s", searchProp, filePath, strSearchRes)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory %s: %v\n", dirPath, err)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run stringgrep.go <path> <uicobjsearch> <objpropsearch>")
		return
	}

	dirPath := os.Args[1]
	searchProp := os.Args[2]
	searchObj := os.Args[3]

	searchInDirectory(dirPath, searchProp, searchObj)
}
