package helper

import (
	"os"
	"log"
)

type ErrorCode int

const (
	SUCCESS = 0
	EOPEN = 1
	EPHOFF = 2
)

func OpenFile(name string) *os.File {
	f, err := os.Open(name)
	if err == nil{
		return f
	}

	if !os.IsNotExist(err) {
		log.Printf("Unexpected error: \"%s\"")
	} else {
		log.Printf("Expected input \"%s\" does not exist.")
	}

	return nil
}



func ReadNextBytesFromFile(input *os.File, count int) []byte{
	data := make([]byte, count)
	_, err := input.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}