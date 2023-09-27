package main

import (
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	"os"
	"strconv"
)

func decodePngFile(fileName string) (int, string) {
	// open and decode image file

	// fmt.Println(arg)
	file, _ := os.Open(fileName)

	defer file.Close()

	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	if err == nil {
		data := result.String()

		n := len(data)

		if data[n-1] == '#' {
			stringChunk := data[:n-1]

			if len(stringChunk) > 4 {
				i, err := strconv.Atoi(stringChunk[:4])
				if err == nil {
					if _, found := fileMap[i]; !found {
						fileMap[i] = fileName
					}
					return i, stringChunk[4:]
				}
			}
		}
		return -1, ``
	} else {
		return -1, ``
	}

}
