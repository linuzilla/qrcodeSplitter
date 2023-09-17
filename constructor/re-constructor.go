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

func reConstructor(comPort string) {
	c := &serial.Config{Name: comPort, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	expectedChunk := 0
	buf := make([]byte, 4096)
	var fileHeader models.FileHeader
	stringChunk := ``

	fileHeader.Count = -1
	var targetFile *os.File

	for {
		if stringChunk == `` {
			if fileHeader.Count == -1 {
				fmt.Println("Please scan the first (qrcode0000.png) file")
			} else {
				fmt.Printf("%d/%d - Please scan qrcode%04d.png file\n", expectedChunk, fileHeader.Count, expectedChunk)
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
						fmt.Printf("please scan #%d QR Code\n", expectedChunk)
					} else {
						decodeString, err := base64.StdEncoding.DecodeString(stringChunk[4:])

						if err != nil {
							fmt.Println(err)
						} else {
							if i == 0 {
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
								}
							}
						}
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
