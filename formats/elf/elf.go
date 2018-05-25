package elf

import (
	"github.com/pspace/golf/helper"
	"fmt"
	"os"
)


func ParserProgramHeaders(f *os.File, count uint16, entrySize uint16) []*ProgramHeaderData{
	phEntries := []*ProgramHeaderData{}
	for i := uint16(0); i < count; i++ {
		ph := ParseProgramHeader(f)
		phEntries = append(phEntries, ph)
	}

	return phEntries
}

func ParseElf(path string){
	f := helper.OpenFile("/bin/ls")
	if nil == f {
		fmt.Printf("An error occurred. Exiting")
		os.Exit(helper.EOPEN)
	}

	elfHeader := ParseELFHeader(f)
	fmt.Printf("%$",(*elfHeader).String())

	phOff := (*elfHeader).E_phoff
	phNum := (*elfHeader).E_phnum
	phEntsize := (*elfHeader).E_phentsize

	if phOff != 0x40 {
		fmt.Errorf("Unexpected PH Offset: %#x\n", phOff)
		os.Exit(helper.EPHOFF)
	}

	f.Seek(int64(phOff), 0)
	programHeaders := ParserProgramHeaders(f, phNum, phEntsize)

	for count, h := range programHeaders{
		fmt.Printf("Header %i:\n%s\n", count, (*h).String())
	}

}

