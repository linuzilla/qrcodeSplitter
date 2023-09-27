package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"qrcodeSplitter/models"
	"strconv"

	"github.com/tarm/serial"
)

func readChunk(s *serial.Port, expectedChunk int, totalCount int) string {
	stringChunk := ``
	buf := make([]byte, 4096)

	for {
		if stringChunk == `` {
			if expectedChunk == 0 {
				fmt.Println("Please scan the first (qrcode0000.png) file")
			} else {
				fmt.Printf("%d/%d - Please scan qrcode%04d.png file\n", expectedChunk, totalCount, expectedChunk)
			}
		}

		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		for n > 0 {
			if buf[n-1] == 10 || buf[n-1] == 13 {
				n--
			} else {
				break
			}
		}
		if n == 0 {
			continue
		}

		newInput := string(buf[:n])

		//fmt.Printf("got input (%d) [%v]\n", n, buf[n-1])

		if buf[n-1] == '#' {
			stringChunk += string(buf[:n-1])

			if len(stringChunk) > 4 {
				i, err := strconv.Atoi(stringChunk[:4])

				if err != nil {
					fmt.Println("unknown chunk")
				} else {
					//fmt.Printf("Chunk #%d\n", i)

					if i != expectedChunk {
						fmt.Printf("please scan #%d QR Code (got %d)\n", expectedChunk, i)
					} else {
						return stringChunk[4:]
					}
				}
			} else {
				fmt.Println("unknown chunk")
			}
			stringChunk = ``
		} else {
			stringChunk += newInput
		}
	}
}

func readChunkFromPngFile(expectedChunk int) string {
	if fileName, found := fileMap[expectedChunk]; found {
		if i, data := decodePngFile(fileName); i == expectedChunk {
			return data
		}
	}
	return ``
}

func reConstructor(comPort string) {
	c := &serial.Config{Name: comPort, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	expectedChunk := 0
	var fileHeader models.FileHeader

	fileHeader.Count = 0
	var targetFile *os.File
	var stringChunk string
	skipPngFile := false

	for {
		if !skipPngFile {
			stringChunk = readChunkFromPngFile(expectedChunk)
			skipPngFile = true
		}

		if stringChunk == `` {
			stringChunk = readChunk(s, expectedChunk, fileHeader.Count)
			skipPngFile = false
		}

		decodeString, err := base64.StdEncoding.DecodeString(stringChunk)

		if err != nil {
			fmt.Println(err)
		} else {
			if expectedChunk == 0 {
				if err != nil {
					fmt.Println(err)
				} else {
					err := json.Unmarshal(decodeString, &fileHeader)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Printf("Filename: %s\n", fileHeader.Filename)
						fmt.Printf("Count: %d\n", fileHeader.Count)

						targetFile, err = os.OpenFile("Copy of - "+fileHeader.Filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
						if err != nil {
							log.Fatal(err)
						}
						expectedChunk++
						skipPngFile = false
					}
				}
			} else {
				_, err = targetFile.Write(decodeString)
				if err != nil {
					fmt.Println(err)
				} else {
					if expectedChunk == fileHeader.Count {
						targetFile.Close()
						fmt.Println("File saved")
						return
					}
					expectedChunk++
					skipPngFile = false
				}
			}
		}
	}
}
