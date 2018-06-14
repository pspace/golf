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
	ECAPSTONE = 3
	EDISASSEMBLE = 4
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



func ReadNextBytesFromFile(input *os.File, count uint64) []byte{
	data := make([]byte, count)
	_, err := input.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func Clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}