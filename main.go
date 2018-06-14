package main

import (
	"github.com/pspace/golf/formats/elf"
	"github.com/pspace/golf/repl"
)

func main() {

	//elf.PrintHeaderDescription()
	//elf.GenerateHeaderStruct()
	//elf.GenerateAll()
	elf.ParseElf("/bin/ls")

	repl.StartInteractive()

	//parsedHeader := elf.ParseELFHeader(f)

	//fmt.Printf("%s", parsedHeader.String())
}
