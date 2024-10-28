package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetClosingObjTag(t *testing.T) {
	strObj := `          <xref object_id="2" screen="pp_area_scroll_bar">`

	strRes := getClosingObjTag(strObj)
	strWant := `</xref>`

	if strRes != strWant {
		t.Fatalf("closing object is not same want:%s get:%s", strWant, strRes)
	}
}

func TestSearchInFileNegative(t *testing.T) {
	filePath := `./dummyfile/negativeres.txt`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`

	obj := SearchCreate(``, ``, ``, searchObj, searchProp)
	res, err := obj.searchInFile(filePath)

	if !(res == ``) || !(err == nil) {
		t.Fatalf(`False, no error or result is  `)
	}
}

func TestSearchInFilePos(t *testing.T) {
	filePath := `./dummyfile/positiveres.txt`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`

	obj := SearchCreate(``, ``, ``, searchObj, searchProp)
	res, err := obj.searchInFile(filePath)

	if (res == ``) || (err != nil) {
		t.Fatalf(`False, no error or result is found `)
	}

	expRes := `Line: 10,`
	if res != expRes {
		t.Fatalf("Wrong location, want: %s get: %s", expRes, res)
	}
}

func TestSearchInFileMixed(t *testing.T) {
	filePath := `./dummyfile/mixedres.txt`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`

	obj := SearchCreate(``, ``, ``, searchObj, searchProp)
	res, err := obj.searchInFile(filePath)

	if (res == ``) || (err != nil) {
		t.Fatalf(`False, no error or result is found `)
	}

	expRes := `Line: 10,`
	if res != expRes {
		t.Fatalf("Wrong location, want: %s get: %s", expRes, res)
	}
}

func TestSearchInFilenotExists(t *testing.T) {
	filePath := `./dummyfile/dummmy.txt`
	searchProp := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchObj := `<property name="type" value="4097"/>`

	obj := SearchCreate(``, ``, ``, searchObj, searchProp)
	_, err := obj.searchInFile(filePath)

	if err == nil {
		t.Fatalf(`False, find shouldnt be found `)
	}
}

func TestGetFileName(t *testing.T) {
	filePathWrld := `C:\Users\iei20120005\Documents\Work\Cate\grepstring_dummy\01_Projects\Gemini_Triforce\Fax\screen\UIWID_FAX_SEND_TOP.nui`
	filePathRel := `./dummyfile/positiveres.txt`

	strResWrld := getFileName(filePathWrld)
	strWantWrld := `UIWID_FAX_SEND_TOP.nui`

	strResRel := getFileName(filePathRel)
	strWantRel := `positiveres.txt`

	if (strResRel != strWantRel) || (strResWrld != strResWrld) {
		if strResRel != strWantRel {
			t.Fatalf(`Filename is not same, want: %s actual:%s`, strWantRel, strResRel)
		} else {
			t.Fatalf(`Filename is not same, want: %s actual:%s`, strWantWrld, strResWrld)
		}
	}
}

func TestMain(t *testing.T) {
	dirPath := `D:\SVN\UIRoot\UIC\01_Projects\Gemini_Triforce`
	searchProp := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchObj := `<property name="type" value="4097"/>`
	fileNameRes := `test_res`
	fileType := `.nui`

	pWorkerPool := createWorkerPool(4)
	pstrGrep := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)

	pWorkerPool.startWorkerPool(pstrGrep)
}

func TestMain2(t *testing.T) {
	dirPath := `D:\SVN\UIRoot\UIC\01_Projects\Gemini_Triforce`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`
	fileNameRes := `test_res2`
	fileType := `.nui`

	sg := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)
	sg.getFilePth()
	for i := 0; i < len(sg.filePth); i++ {
		fmt.Printf("searching in:%s\n", sg.filePth[i])
		res, _ := sg.searchInFile(sg.filePth[i])
		if res != `` {
			fmt.Printf("writing:%s\n", res)
			sg.writeResCSV(res)
		}
	}
}

func TestWorkerPool(t *testing.T) {
	dirPath := `D:\SVN\UIRoot\UIC\01_Projects\Gemini_Triforce`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`
	fileNameRes := `test_res2`
	fileType := `.nui`

	sg := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)
	wp := createWorkerPool(4)
	wp.startWorkerPool(sg)
}

func TestConccurentWalk(t *testing.T) {
	dirPath := `D:\SVN\UIRoot\UIC\01_Projects\Gemini_Triforce`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`
	fileNameRes := `test_res3`
	fileType := `.nui`

	sg := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(*sync.WaitGroup) {
			defer wg.Done()
			sg.getFilePth()
		}(&wg)
	}

	wg.Wait()

	for i := 0; i < len(sg.filePth); i++ {
		fmt.Println(sg.filePth[i])
	}

}

func TestWorkerJob(t *testing.T) {
	dirPath := `D:\SVN\UIRoot\UIC\01_Projects\Gemini_Triforce`
	searchObj := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchProp := `<property name="type" value="3072"/>`
	fileNameRes := ``
	fileType := `.nui`

	sg := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)
	sg2 := SearchCreate(dirPath, fileType, fileNameRes, searchObj, searchProp)

	chDir := make(chan string)
	chWalk := make(chan string)
	chErr := make(chan error)
	//chWrite := make(chan string)

	var wg sync.WaitGroup
	var pFlag *bool
	bFlag := false
	pFlag = &bFlag

	wGet := newWorker(1)
	wg.Add(1)

	go wGet.WalkJobGetter(sg, &wg)

	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := range sg.dirFndPath {
			chDir <- sg.dirFndPath[i]
		}

		chDir <- `done`

	}()

	for i := 0; i < 2; i++ {
		w := newWorker(1)
		wg.Add(1)
		go w.WorkerJob(chWalk, chDir, jobWalk, 0, 0, sg, &wg, pFlag, chErr)
	}
	wg.Wait()

	close(chWalk)

	sg2.getFilePth()

	if len(sg.filePth) != len(sg2.filePth) {
		t.Fatalf("length res is not same %d|%d\n", len(sg.filePth), len(sg2.filePth))
	}
}
