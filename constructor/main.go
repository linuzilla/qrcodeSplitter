package main

import "flag"

func main() {
	comPort := flag.String("com", "COM3", "com port")
	flag.Parse()

	reConstructor(*comPort)
}
