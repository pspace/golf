package elf

import (
	"encoding/binary"
	"github.com/pspace/golf/helper"
	"os"
)

var ELFByteOrder binary.ByteOrder = binary.LittleEndian

func setByteOrder(indicator byte) {
	if indicator == 1 {
		ELFByteOrder = binary.LittleEndian
	}
	if indicator == 2 {
		ELFByteOrder = binary.BigEndian
	}
}

func ParseELFHeader(input *os.File) *ELFHeaderData {
	elfHeaderData := ELFHeaderData{}
	var data []byte
	data = helper.ReadNextBytesFromFile(input, 4)
	elfHeaderData.EI_MAG = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 1)
	elfHeaderData.EI_CLASS = data[0]

	data = helper.ReadNextBytesFromFile(input, 1)
	elfHeaderData.EI_DATA = data[0]
	setByteOrder(data[0])

	data = helper.ReadNextBytesFromFile(input, 1)
	elfHeaderData.EI_VERSION = data[0]

	data = helper.ReadNextBytesFromFile(input, 1)
	elfHeaderData.EI_OSABI = data[0]

	data = helper.ReadNextBytesFromFile(input, 8)
	elfHeaderData.EI_PAD = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_type = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_machine = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	elfHeaderData.E_version = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	elfHeaderData.E_entry = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	elfHeaderData.E_phoff = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	elfHeaderData.E_shoff = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	elfHeaderData.E_flags = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_ehsize = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_phentsize = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_phnum = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_shentsize = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_shnum = ELFByteOrder.Uint16(data)

	data = helper.ReadNextBytesFromFile(input, 2)
	elfHeaderData.E_shstrndx = ELFByteOrder.Uint16(data)

	return &elfHeaderData
}

func ParseProgramHeader(input *os.File) *ProgramHeaderData {
	headerData := ProgramHeaderData{}
	var data []byte
	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.P_type = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.P_flags = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_offset = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_vaddr = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_addr = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_filesz = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_memsz = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.P_align = ELFByteOrder.Uint64(data)

	return &headerData
}

func ParseSectionHeader(input *os.File) *SectionHeaderData {
	headerData := SectionHeaderData{}
	var data []byte
	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.SH_name = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.SH_type = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_flags = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_addr = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_offset = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_size = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.SH_link = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 4)
	headerData.SH_info = ELFByteOrder.Uint32(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_addralign = ELFByteOrder.Uint64(data)

	data = helper.ReadNextBytesFromFile(input, 8)
	headerData.SH_entsize = ELFByteOrder.Uint64(data)

	return &headerData
}
