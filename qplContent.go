package qpl

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"strings"
)

type Identifier struct {
	Catalog string `xml:"Catalog,attr"`
	Entry   string `xml:"Entry,attr"`
}

type LocalizableText struct {
	Language string `xml:"Language,attr"`
	Text     string `xml:",chardata"`
}

type General struct {
	Structure   string          `xml:"Structure,attr"`
	Identifier  Identifier      `xml:"Identifier"`
	Title       LocalizableText `xml:"Title"`
	Language    LocalizableText `xml:"Language"`
	Description LocalizableText `xml:"Description"`
	Keyword     LocalizableText `xml:"Keyword"`
}

type MetaData struct {
	General General `xml:"General"`
}

type QPLContentObject struct {
	ContentType string   `xml:"Type,attr"`
	MetaData    MetaData `xml:"MetaData"`
}

func readQPLXMLFile(fileEntry fs.DirEntry) (QPLContentObject, error) {
	file, err := os.Open(fileEntry.Name())
	if err != nil {
		return QPLContentObject{}, err
	}

	defer file.Close()

	decoder := xml.NewDecoder(file)

	var qplContentObject QPLContentObject
	err = decoder.Decode(&qplContentObject)
	if err != nil {
		return QPLContentObject{}, err
	}

	return qplContentObject, nil
}

func findQTIFile(contents []os.DirEntry) (os.DirEntry, error) {
	for _, content := range contents {
		if content.IsDir() {
			continue
		}

		if strings.Contains(content.Name(), "qti") && strings.HasSuffix(content.Name(), ".xml") {
			return content, nil
		}
	}

	return nil, errors.New("qti file not found")
}

func findQPLFile(contents []os.DirEntry) (os.DirEntry, error) {
	for _, content := range contents {
		if content.IsDir() {
			continue
		}

		if strings.Contains(content.Name(), "qpl") && strings.HasSuffix(content.Name(), ".xml") {
			return content, nil
		}
	}

	return nil, errors.New("qpl file not found")
}

func ReadQPLFolder(folder string) (QPLFile, error) {
	// check if folder exists and is a directory
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return QPLFile{}, err
	}

	// read all files and folder in folder
	contents, err := os.ReadDir(folder)
	if err != nil {
		return QPLFile{}, err
	}

	// if contents count is 1, then make contents = os.ReadDir of contents[0]
	if len(contents) == 1 {
		contents, err = os.ReadDir(folder + "/" + contents[0].Name())
		if err != nil {
			return QPLFile{}, err
		}
	}

	// find qpl file (file which contains "qpl" in the name and has xml file extension)
	qplFile, err := findQPLFile(contents)
	if err != nil {
		return QPLFile{}, err
	}

	// qti file
	qtiFile, err := findQTIFile(contents)
	if err != nil {
		return QPLFile{}, err
	}

	// read qplFile and qtiFile
	qplContentObject, err := readQPLXMLFile(qplFile)
	if err != nil {
		return QPLFile{}, err
	}

	qtiContentObject, err := ReadQTIXMLFile(qtiFile)
	if err != nil {
		return QPLFile{}, err
	}

	return QPLFile{
		ContentObject:   qplContentObject,
		QuestestInterop: qtiContentObject,
	}, nil
}
