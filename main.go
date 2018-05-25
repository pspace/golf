package main

import "github.com/pspace/golf/formats/elf"

func main() {

	//elf.PrintHeaderDescription()
	//elf.GenerateHeaderStruct()

	elf.ParseElf("/bin/ls")

	//parsedHeader := elf.ParseELFHeader(f)

	//fmt.Printf("%s", parsedHeader.String())
}
