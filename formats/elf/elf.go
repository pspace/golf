package elf

import (
	"github.com/pspace/golf/helper"
	"fmt"
	"os"
)

func ParseProgramHeaders(f *os.File, count uint16, entrySize uint16) []*ProgramHeaderData{
	phEntries := []*ProgramHeaderData{}
	for i := uint16(0); i < count; i++ {
		ph := ParseProgramHeader(f)
		phEntries = append(phEntries, ph)
	}

	return phEntries
}


func ParseSectionHeaders(f *os.File, count uint16, entrySize uint16) []*SectionHeaderData{
	shEntries := []*SectionHeaderData{}
	for i := uint16(0); i < count; i++ {
		sh := ParseSectionHeader(f)
		shEntries = append(shEntries, sh)
	}

	return shEntries
}

func GetProgram(f *os.File, elfHeader *ELFHeaderData)  {

	phOff := (*elfHeader).E_phoff
	phNum := (*elfHeader).E_phnum
	phEntsize := (*elfHeader).E_phentsize

	if phOff != 0x40 {
		fmt.Errorf("Unexpected PH Offset: %#x\n", phOff)
		os.Exit(helper.EPHOFF)
	}

	f.Seek(int64(phOff), 0)
	programHeaders := ParseProgramHeaders(f, phNum, phEntsize)

	for count, h := range programHeaders{
		fmt.Printf("Header %d:\n%s\n", count, (*h).String())
	}
}

func GetSections(f *os.File, elfHeader *ELFHeaderData) {
	shOff := (*elfHeader).E_shoff
	shNum := (*elfHeader).E_shnum
	shEntsize := (*elfHeader).E_shentsize

	f.Seek(int64(shOff), 0)
	var sectionNames []byte

	sectionHeaders := ParseSectionHeaders(f, shNum, shEntsize)
	nameHeader := sectionHeaders[(*elfHeader).E_shstrndx]
	f.Seek(int64((*nameHeader).SH_offset), 0)
	sectionNames = helper.ReadNextBytesFromFile(f, (*nameHeader).SH_size)

	for count, h := range sectionHeaders{
		nameOffset := (*h).SH_name
		nameSize := helper.Clen(sectionNames[(*h).SH_name:])

		name := string(sectionNames[nameOffset : int(nameOffset) + nameSize])
		fmt.Printf("%d. Section %s - %s:\n%s\n", count, name, SH_name_map[name], (*h).String())

	}

}


func ParseElf(path string){
	f := helper.OpenFile("/bin/ls")
	if nil == f {
		fmt.Printf("An error occurred. Exiting")
		os.Exit(helper.EOPEN)
	}

	elfHeader := ParseELFHeader(f)
	fmt.Printf("%$",(*elfHeader).String())

	//GetProgram(f, elfHeader)
	GetSections(f, elfHeader)


}

