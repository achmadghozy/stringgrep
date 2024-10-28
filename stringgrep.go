package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Common func library

func getClosingObjTag(stringTag string) string {

	strRaw := strings.Fields(stringTag)
	strRes := strings.Replace(strRaw[0], `<`, ``, -1)

	return (`</` + strRes + `>`)
}

func getFileName(strFilePth string) string {
	var tempArr []string
	if strings.Contains(strFilePth, `:\`) {
		tempArr = strings.Split(strFilePth, `\`)
	} else {
		tempArr = strings.Split(strFilePth, `/`)
	}
	return tempArr[len(tempArr)-1]
}

func writeToCSV(strResFilePth string, strFindRes string) error {
	f, err := os.OpenFile(strResFilePth, os.O_APPEND, 0)

	if err != nil {
		fmt.Printf("Error opening the file\n")
		return err
	}
	defer f.Close()

	_, err = fmt.Fprint(f, strFindRes)

	/*
		_, err = f.WriteString(strFindRes)
	*/

	if err != nil {
		fmt.Printf("Error writing to result\n")
		return err
	}

	err = f.Sync()
	if err != nil {
		fmt.Printf("Error saving the result\n")
		return err
	}

	return nil
}

type stUICObjGrep struct {
	bFirstExc   bool     // is First execution
	dirPath     string   // Directory path where uic proj is located
	fileType    string   // Filetype of target
	fileNameRes string   // Filename where result is saved
	searchObj   string   // Object tag to be searched
	searchProp  string   // Object property to be searched
	filReesPth  string   // File Res path
	filePth     []string // Filepath where dir is being walked
	dirFndPath  []string // Directories where dir is being walked
}

func SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp string) *stUICObjGrep {
	filReesPth := `./result/` + fileNameRes + `.csv`
	bFirstExc := true

	return &stUICObjGrep{
		dirPath:     dirPath,
		fileType:    fileType,
		fileNameRes: fileNameRes,
		searchObj:   searchObj,
		searchProp:  searchProp,
		filReesPth:  filReesPth,
		bFirstExc:   bFirstExc,
	}
}

func (og *stUICObjGrep) createResCSV() error {
	_, err := os.Create(og.filReesPth)

	if err != nil {
		fmt.Printf("File cannot be created\n")
		return err
	}

	return nil
}

func (og *stUICObjGrep) writeResCSV(strRes string) error {
	var err error

	if og.bFirstExc {
		_, err := os.Stat(og.filReesPth)
		if errors.Is(err, os.ErrNotExist) {
			og.createResCSV()
			og.bFirstExc = false
		} else {
			fmt.Printf("Error creating the file\n")
			return err
		}
	}

	if err != nil {
		return err
	}

	err = writeToCSV(og.filReesPth, strRes)
	if err != nil {
		return err
	}

	return nil
}

/*****************************************************************************************/
/*
Function returns:

	Strings		: if found return FilePth, Line Obj, Line Obj Prop, if found more than one append newlines
	err			: return opening function if fail to open the file

Function args:

	filePath	: file path
	searchProp	: property tag to search
	searchObj	: object tag to search
*/
/*****************************************************************************************/
func (og *stUICObjGrep) searchInFile(filePath string) (string, error) {
	var strObjTag string
	var strBuilderObj strings.Builder
	var strBuilderObj2 strings.Builder

	line := 1
	bCandFound := false
	// filePathRes := `./result/` + filenameRes + `.csv` /* Write to CSV is moved to main func */
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s\n", filePath)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// remove empty space fields
		strTempObj := strings.Join(strings.Fields(scanner.Text()), ` `)
		switch bCandFound {

		// If candidate is found
		case true:
			{
				if (strings.Contains(strTempObj, og.searchProp)) || (strings.Contains(strTempObj, strObjTag)) {
					bCandFound = false
					// If property was found inside the object tag, or the object tag ends
					// If property was found inside, assign result
					if strings.Contains(strTempObj, og.searchProp) {

						// If found write the line number of obj props
						strBuilderObj.WriteString(strconv.Itoa(line))
						/*	Write to CSV is moved to main func
							if filenameRes != `` {
								strWriteToCSV := fmt.Sprintf(`%s,%s`, getFileName(filePath), strconv.Itoa(line))
								err = writeToCSV(filePathRes, strWriteToCSV)
								if err != nil {
									return ``, err
								}
							}
						*/
					}

					// If obj tag ends write the result
					strBuilderObj2.WriteString(fmt.Sprintln(strBuilderObj.String()))
				}
				break
			}

		// Default search for object
		default:
			{
				if strings.Contains(strTempObj, og.searchObj) {

					// If found reset the content of the string and write the line number of obj
					strBuilderObj.Reset()
					strBuilderObj.WriteString(filePath + `,` + strconv.Itoa(line) + `,`)

					// Check if the object has closing tags, therefore properties are not overriden
					if strings.HasPrefix(strTempObj, "<") && strings.HasSuffix(strTempObj, "/>") {
						// do nothing
					} else {
						bCandFound = true
						strObjTag = getClosingObjTag(strTempObj)
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

	strResF := strBuilderObj2.String()

	return strResF, nil
}

/*****************************************************************************************/
/*
Function returns:

	[]strings	: return array of filepath
	err			: return opening function if fail to open the file

Function args:

	filePath	: Dir/ File path
	fileNameRes	: Result written in named file .csv
	fileType	: file type where object want to be searched
*/
/*****************************************************************************************/
func (og *stUICObjGrep) getFilePth() ([]string, error) {

	if og.fileType == `` {
		og.fileType = `.`
	}

	err := filepath.WalkDir(og.dirPath, func(filePath string, info os.DirEntry, err error) error {
		// fmt.Printf("Walking %s\n", filePath)
		if err != nil {
			fmt.Printf("Error walking the path %s: %v\n", filePath, err)
			return err
		}
		if !info.IsDir() {
			if strings.Contains(getFileName(filePath), og.fileType) {

				// fmt.Printf(`walking path: %s \n`, filePath)
				// Save the filepath to be used on main
				og.filePth = append(og.filePth, filePath)

			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory %s: %v\n", og.dirPath, err)
	}

	return og.filePth, nil
}

// Get all directories should be invoked before getfilePth2
func (og *stUICObjGrep) getAllDirectories() ([]string, error) {
	entries, err := os.ReadDir(og.dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {

			// if folder is hidden don't add
			if !strings.HasPrefix(entry.Name(), `.`) {
				og.dirFndPath = append(og.dirFndPath, fmt.Sprintf("%s\\%s", og.dirPath, entry.Name()))
			}
		}
	}

	return og.dirFndPath, nil
}

// Should be invoked after GetAllDirectories
func (og *stUICObjGrep) getFilePth2(dirPath string) ([]string, error) {
	if og.fileType == `` {
		og.fileType = `.`
	}

	err := filepath.WalkDir(dirPath, func(filePath string, info os.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("Error walking the path %s: %v", filePath, err)

		}
		if !info.IsDir() {
			if strings.Contains(getFileName(filePath), og.fileType) {

				// append the list of filepath in class variable
				og.filePth = append(og.filePth, filePath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error walking the directory %s: %v", dirPath, err)
	}

	return og.filePth, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run stringgrep.go <path> <uicobjsearch> <objpropsearch> Opt:<filename_res> <filetypes>")
		return
	}
	dirPath := os.Args[1]
	searchProp := os.Args[2]
	searchObj := os.Args[3]
	filenameRes := os.Args[4]
	fileType := os.Args[5]

	pStrGrep := SearchCreate(dirPath, searchProp, searchObj, filenameRes, fileType)
	_, err := pStrGrep.getFilePth()
	if err != nil {
		fmt.Printf("Error in walking path\n")
	}

	/*
		pStrGrep2 := SearchCreate(dirPath, searchProp, searchObj, filenameRes, fileType)
		pWP := createWorkerPool(4)
		pWP.startWorkerPool(pStrGrep2)
	*/
}
