package models

type FileHeader struct {
	Filename string `json:"filename"`
	Count    int    `json:"count"`
}
