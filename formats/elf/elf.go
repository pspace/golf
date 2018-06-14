package elf

import (
	"github.com/pspace/golf/helper"
	"fmt"
	"os"
	"github.com/pspace/golf/disassembler"
	"github.com/bnagy/gapstone"
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

func GetSections(f *os.File, elfHeader *ELFHeaderData) []*SectionHeaderData{
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

	return sectionHeaders
}


func ParseElf(path string){
	f := helper.OpenFile(path)
	if nil == f {
		fmt.Printf("An error occurred. Exiting")
		os.Exit(helper.EOPEN)
	}

	elfHeader := ParseELFHeader(f)
	fmt.Printf("%s",(*elfHeader).String())

	//GetProgram(f, elfHeader)
	sectionHeaders := GetSections(f, elfHeader)

	for _, h := range sectionHeaders{
		t := (*h).SH_type
		if t == SECTION_TYPE_PROGBITS {
			disasmInfo := disassembler.DisassemblerSpecs{gapstone.CS_ARCH_X86, gapstone.CS_MODE_64}
			f.Seek( int64((*h).SH_offset), 0)
			data := helper.ReadNextBytesFromFile(f,  (*h).SH_size)
			disassembler.Disassemble(disasmInfo, &data)
		}
	}


}