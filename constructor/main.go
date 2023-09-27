package main

import (
	"flag"
)

var fileMap map[int]string

func main() {
	comPort := flag.String("com", "COM3", "com port")
	flag.Parse()

	readPngFiles()

	reConstructor(*comPort)
}
