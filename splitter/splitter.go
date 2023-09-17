package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"qrcodeSplitter/models"
)

func generateQrCode(count int, encodedData string, qrCodeDir string, imageSize int) (string, error) {
	qrCode, err := qrcode.New(encodedData, qrcode.Low)

	if err != nil {
		return "", err
	}

	// Create a QR code image file name
	qrCodeFileName := fmt.Sprintf("%s/qrcode%04d.png", qrCodeDir, count)

	// Create a QR code image file
	qrCodeFile, err := os.Create(qrCodeFileName)
	if err != nil {
		return "", err
	}
	defer qrCodeFile.Close()

	// Save the QR code image as PNG
	err = png.Encode(qrCodeFile, qrCode.Image(imageSize))

	return qrCodeFileName, err
}

func splitter(inputFileName string, imageSize int, chunkSize int) {
	// Specify the input file name

	// Open the input file for reading
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Unable to open the input file:", err)
		return
	}
	defer inputFile.Close()

	// Create a buffer to store file chunks
	buffer := make([]byte, chunkSize)

	// Create a directory to store QR code images
	qrCodeDir := "qrcodes"
	if _, err := os.Stat(qrCodeDir); os.IsNotExist(err) {
		os.Mkdir(qrCodeDir, os.ModeDir|os.ModePerm)
	}

	// err = os.MkdirAll(qrCodeDir, os.ModePerm) // Ensure directory exists with proper permissions

	// Counter for generating QR code file names
	count := 1

	for {
		n, err := inputFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading the file:", err)
			return
		}

		// Encode the read chunk into Base64
		encodedData := fmt.Sprintf("%04d", count) + base64.StdEncoding.EncodeToString(buffer[:n]) + "#"

		qrCodeFileName, err := generateQrCode(count, encodedData, qrCodeDir, imageSize)
		// Generate a QR code image
		if err != nil {
			fmt.Println("Error generating QR code:", err)
			return
		}

		// Create a QR code image file name
		fmt.Printf("Generated QR code file: %s\n", qrCodeFileName)
		count++
	}

	marshal, err := json.Marshal(&models.FileHeader{
		Filename: filepath.Base(inputFileName),
		Count:    count - 1,
	})

	encodedData := fmt.Sprintf("%04d", 0) + base64.StdEncoding.EncodeToString(marshal) + "#"
	_, err = generateQrCode(0, encodedData, qrCodeDir, imageSize)
	// Generate a QR code image
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println(string(marshal))
}
