package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func readPngFiles() {
	// open and decode image file
	var pngFiles []string

	fileMap = make(map[int]string)

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".PNG")) {
			if match, _ := regexp.MatchString(`^.*[0-9]{4}\.`, file.Name()); match {
				pngFiles = append(pngFiles, file.Name())
			}
		}
	}

	sort.Strings(pngFiles)

	for _, fileName := range pngFiles {
		if i, data := decodePngFile(fileName); i >= 0 && len(data) > 4 {
			if _, found := fileMap[i]; !found {
				fileMap[i] = fileName
			}
		}
	}
	fmt.Println(fileMap)
}
