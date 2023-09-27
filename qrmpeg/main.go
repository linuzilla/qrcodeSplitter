package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	// open and decode image file
	var pngFiles []string

	prev := ``
	fileMap := make(map[int]string)

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

	for _, arg := range pngFiles {
		// fmt.Println(arg)
		file, _ := os.Open(arg)
		img, _, _ := image.Decode(file)

		// prepare BinaryBitmap
		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

		// decode image
		qrReader := qrcode.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)

		if err == nil {
			data := result.String()

			if prev != data {
				// fmt.Println(result)
				n := len(data)

				if data[n-1] == '#' {
					stringChunk := data[:n-1]

					if len(stringChunk) > 4 {
						i, err := strconv.Atoi(stringChunk[:4])
						if err == nil {
							if _, found := fileMap[i]; !found {
								fileMap[i] = arg
							}
						}
					}
				}
			}
			prev = data
		}
	}
	fmt.Println(fileMap)
}
