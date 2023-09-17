package main

import "flag"

func main() {
	filename := flag.String("file", "", "encode file")
	chunkSize := flag.Int("chunk", 768, "chunk chunkSize")
	imageSize := flag.Int("size", 384, "size imageSize")
	flag.Parse()

	if *filename != "" {
		splitter(*filename, *imageSize, *chunkSize)
	} else {
		flag.Usage()
	}
}
