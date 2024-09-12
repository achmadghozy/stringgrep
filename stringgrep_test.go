package main

import (
	"testing"
)

func TestSearchInFileNegative(t *testing.T) {
	filePath := `./dummyfile/negativeres.txt`
	searchProp := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchObj := `<property name="type" value="4097"/>`

	res, err := searchInFile(filePath, searchProp, searchObj)

	if !(res == ``) || !(err == nil) {
		t.Fatalf(`False, no error or result is  `)
	}
}

func TestSearchInFilePos(t *testing.T) {
	filePath := `./dummyfile/positiveres.txt`
	searchProp := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchObj := `<property name="type" value="4097"/>`

	res, err := searchInFile(filePath, searchProp, searchObj)

	if !(res == ``) || !(err == nil) {
		t.Fatalf(`False, no error or result is found `)
	}
}

func TestSearchInFilenotExists(t *testing.T) {
	filePath := `./dummyfile/dummmy.txt`
	searchProp := `<xref object_id="2" screen="pp_area_scroll_bar">`
	searchObj := `<property name="type" value="4097"/>`

	_, err := searchInFile(filePath, searchProp, searchObj)

	if err != nil {
		t.Fatalf(`False, find shouldnt be found `)
	}
}
