package elf

import (
	"fmt"
	"strings"
)

type HeaderDescriptionEntry struct {
	offset      uint
	size        uint
	field       string
	description string
}

//https://www.freebsd.org/cgi/man.cgi?query=elf&sektion=5
var E_OSABI_map = map[byte]string{
	0:   "UNIX System V ABI",
	1:   "HP-UX operating system",
	2:   "NetBSD",
	3:   "GNU/Linux",
	4:   "GNU/Hurd",
	6:   "Solaris",
	7:   "AIX",
	8:   "IRIX",
	9:   "FreeBSD",
	10:  "TRU64 UNIX",
	11:  "Novell Modesto",
	12:  "OpenBSD",
	13:  "OpenVMS",
	14:  "Hewlett-Packard Non-Stop Kernel",
	15:  "AROS",
	16:  "FenixOS",
	17:  "Nuxi CloudABI",
	64:  "First architecture-specific OS ABI/AMD HSA runtime/Bare-metal TMS320C6000",
	65:  "AMD PAL runtime/Linux TMS320C6000",
	66:  "AMD GCN GPUs (GFX6+) for MESA runtime",
	97:  "ARM",
	255: "Standalone (embedded) application",
};

//http://llvm.org/doxygen/Support_2ELF_8h_source.html
var E_machine_map = map[uint16]string{
	0:      "No machine",
	1:      "AT&T WE 32100",
	2:      "SPARC",
	3:      "Intel 386",
	4:      "Motorola 68000",
	5:      "Motorola 88000",
	6:      "Intel MCU",
	7:      "Intel 80860",
	8:      "MIPS R3000",
	9:      "IBM System/370",
	10:     "MIPS RS3000 Little-endian",
	15:     "Hewlett-Packard PA-RISC",
	17:     "Fujitsu VPP500",
	18:     "Enhanced instruction set SPARC",
	19:     "Intel 80960",
	20:     "PowerPC",
	21:     "PowerPC64",
	22:     "IBM System/390",
	23:     "IBM SPU/SPC",
	36:     "NEC V800",
	37:     "Fujitsu FR20",
	38:     "TRW RH-32",
	39:     "Motorola RCE",
	40:     "ARM",
	41:     "DEC Alpha",
	42:     "Hitachi SH",
	43:     "SPARC V9",
	44:     "Siemens TriCore",
	45:     "Argonaut RISC Core",
	46:     "Hitachi H8/300",
	47:     "Hitachi H8/300H",
	48:     "Hitachi H8S",
	49:     "Hitachi H8/500",
	50:     "Intel IA-64 processor architecture",
	51:     "Stanford MIPS-X",
	52:     "Motorola ColdFire",
	53:     "Motorola M68HC12",
	54:     "Fujitsu MMA Multimedia Accelerator",
	55:     "Siemens PCP",
	56:     "Sony nCPU embedded RISC processor",
	57:     "Denso NDR1 microprocessor",
	58:     "Motorola Star*Core processor",
	59:     "Toyota ME16 processor",
	60:     "STMicroelectronics ST100 processor",
	61:     "Advanced Logic Corp. TinyJ embedded processor family",
	62:     "AMD x86-64 architecture",
	63:     "Sony DSP Processor",
	64:     "Digital Equipment Corp. PDP-10",
	65:     "Digital Equipment Corp. PDP-11",
	66:     "Siemens FX66 microcontroller",
	67:     "STMicroelectronics ST9+ 8/16 bit microcontroller",
	68:     "STMicroelectronics ST7 8-bit microcontroller",
	69:     "Motorola MC68HC16 Microcontroller",
	70:     "Motorola MC68HC11 Microcontroller",
	71:     "Motorola MC68HC08 Microcontroller",
	72:     "Motorola MC68HC05 Microcontroller",
	73:     "Silicon Graphics SVx",
	74:     "STMicroelectronics ST19 8-bit microcontroller",
	75:     "Digital VAX",
	76:     "Axis Communications 32-bit embedded processor",
	77:     "Infineon Technologies 32-bit embedded processor",
	78:     "Element 14 64-bit DSP Processor",
	79:     "LSI Logic 16-bit DSP Processor",
	80:     "Donald Knuth's educational 64-bit processor",
	81:     "Harvard University machine-independent object files",
	82:     "SiTera Prism",
	83:     "Atmel AVR 8-bit microcontroller",
	84:     "Fujitsu FR30",
	85:     "Mitsubishi D10V",
	86:     "Mitsubishi D30V",
	87:     "NEC v850",
	88:     "Mitsubishi M32R",
	89:     "Matsushita MN10300",
	90:     "Matsushita MN10200",
	91:     "picoJava",
	92:     "OpenRISC 32-bit embedded processor",
	93:     "ARC International ARCompact processor",
	94:     "Tensilica Xtensa Architecture",
	95:     "Alphamosaic VideoCore processor",
	96:     "Thompson Multimedia General Purpose Processor",
	97:     "National Semiconductor 32000 series",
	98:     "Tenor Network TPC processor",
	99:     "Trebia SNP 1000 processor",
	100:    "STMicroelectronics (www.st.com) ST200",
	101:    "Ubicom IP2xxx microcontroller family",
	102:    "MAX Processor",
	103:    "National Semiconductor CompactRISC microprocessor",
	104:    "Fujitsu F2MC16",
	105:    "Texas Instruments embedded microcontroller msp430",
	106:    "Analog Devices Blackfin (DSP) processor",
	107:    "S1C33 Family of Seiko Epson processors",
	108:    "Sharp embedded microprocessor",
	109:    "Arca RISC Microprocessor",
	110:    "Microprocessor series from PKU-Unity Ltd. and MPRC of Peking University",
	111:    "eXcess: 16/32/64-bit configurable embedded CPU",
	112:    "Icera Semiconductor Inc. Deep Execution Processor",
	113:    "Altera Nios II soft-core processor",
	114:    "National Semiconductor CompactRISC CRX",
	115:    "Motorola XGATE embedded processor",
	116:    "Infineon C16x/XC16x processor",
	117:    "Renesas M16C series microprocessors",
	118:    "Microchip Technology dsPIC30F Digital Signal Controller",
	119:    "Freescale Communication Engine RISC core",
	120:    "Renesas M32C series microprocessors",
	131:    "Altium TSK3000 core",
	132:    "Freescale RS08 embedded processor",
	133:    "Analog Devices SHARC family of 32-bit DSP processors",
	134:    "Cyan Technology eCOG2 microprocessor",
	135:    "Sunplus S+core7 RISC processor",
	136:    "New Japan Radio (NJR) 24-bit DSP Processor",
	137:    "Broadcom VideoCore III processor",
	138:    "RISC processor for Lattice FPGA architecture",
	139:    "Seiko Epson C17 family",
	140:    "The Texas Instruments TMS320C6000 DSP family",
	141:    "The Texas Instruments TMS320C2000 DSP family",
	142:    "The Texas Instruments TMS320C55x DSP family",
	160:    "STMicroelectronics 64bit VLIW Data Signal Processor",
	161:    "Cypress M8C microprocessor",
	162:    "Renesas R32C series microprocessors",
	163:    "NXP Semiconductors TriMedia architecture family",
	164:    "Qualcomm Hexagon processor",
	165:    "Intel 8051 and variants",
	166:    "STMicroelectronics STxP7x family of configurable and extensible RISC processors",
	167:    "Andes Technology compact code size embedded RISC processor family",
	168:    "Cyan Technology eCOG1X family",
	169:    "Dallas Semiconductor MAXQ30 Core Micro-controllers",
	170:    "New Japan Radio (NJR) 16-bit DSP Processor",
	171:    "M2000 Reconfigurable RISC Microprocessor",
	172:    "Cray Inc. NV2 vector architecture",
	173:    "Renesas RX family",
	174:    "Imagination Technologies META processor architecture",
	175:    "MCST Elbrus general purpose hardware architecture",
	176:    "Cyan Technology eCOG16 family",
	177:    "National Semiconductor CompactRISC CR16 16-bit microprocessor",
	178:    "Freescale Extended Time Processing Unit",
	179:    "Infineon Technologies SLE9X core",
	180:    "Intel L10M",
	181:    "Intel K10M",
	183:    "ARM AArch64",
	185:    "Atmel Corporation 32-bit microprocessor family",
	186:    "STMicroeletronics STM8 8-bit microcontroller",
	187:    "Tilera TILE64 multicore architecture family",
	188:    "Tilera TILEPro multicore architecture family",
	190:    "NVIDIA CUDA architecture",
	191:    "Tilera TILE-Gx multicore architecture family",
	192:    "CloudShield architecture family",
	193:    "KIPO-KAIST Core-A 1st generation processor family",
	194:    "KIPO-KAIST Core-A 2nd generation processor family",
	195:    "Synopsys ARCompact V2",
	196:    "Open8 8-bit RISC soft processor core",
	197:    "Renesas RL78 family",
	198:    "Broadcom VideoCore V processor",
	199:    "Renesas 78KOR family",
	200:    "Freescale 56800EX Digital Signal Controller (DSC)",
	201:    "Beyond BA1 CPU architecture",
	202:    "Beyond BA2 CPU architecture",
	203:    "XMOS xCORE processor family",
	204:    "Microchip 8-bit PIC(r) family",
	205:    "Reserved by Intel",
	206:    "Reserved by Intel",
	207:    "Reserved by Intel",
	208:    "Reserved by Intel",
	209:    "Reserved by Intel",
	210:    "KM211 KM32 32-bit processor",
	211:    "KM211 KMX32 32-bit processor",
	212:    "KM211 KMX16 16-bit processor",
	213:    "KM211 KMX8 8-bit processor",
	214:    "KM211 KVARC processor",
	215:    "Paneve CDP architecture family",
	216:    "Cognitive Smart Memory Processor",
	217:    "iCelero CoolEngine",
	218:    "Nanoradio Optimized RISC",
	219:    "CSR Kalimba architecture family",
	224:    "AMD GPU architecture",
	243:    "RISC-V",
	244:    "Lanai 32-bit processor",
	247:    "Linux kernel bpf virtual machine",
	0x4157: "Web Assembly",
};

var ELFHeaderDescription = []HeaderDescriptionEntry{
	{
		0, 4, "EI_MAG", "Magic bytes",
	},
	{
		4, 1, "EI_CLASS", "32 Bit ( == 1) or 64 Bit ( == 2)",
	},
	{
		5, 1, "EI_DATA", "Endianess: 1 - little, 2 - big",
	},
	{
		6, 1, "EI_VERSION", "ELF Version",
	},
	{
		7, 1, "EI_OSABI", "Target plattform",
	},
	{
		8, 0, "EI_ABIVERSION", "ABI Version, since Linux 2.6 gone, EI_PAD grown by 1",
	},
	{
		8, 8, "EI_PAD", "Apple sponsored content",
	},
	{
		0x10, 2, "E_type", "1 relocatable, 2 executable, 3 shared, 4 core",
	},
	{
		0x12, 2, "E_machine",
		"Instruction set architecture, 0x3 x86, 0x28 ARM, 0x3E x64, 0xB7 AArch64, 0xF3 RISC-V",
	},
	{
		0x14, 4, "E_version", "Version - 1 for original ELF",
	},
	{
		0x18, 8, "E_entry", "Entrypoint",
	},
	{
		0x20, 8, "E_phoff", "Programm header table offset (should be 0x40)",
	},
	{
		0x28, 8, "E_shoff", "Start of section header table",
	},
	{
		0x30, 4, "E_flags", "architecture specific flags",
	},
	{
		0x34, 2, "E_ehsize", "Size of the ELF header - should be 64 (0x40) bytes",
	},
	{
		0x36, 2, "E_phentsize", "Size of a program header table entry",
	},
	{
		0x38, 2, "E_phnum", "Number of program header table entries",
	},
	{
		0x3A, 2, "E_shentsize", "Size of a section header table entry",
	},
	{
		0x3C, 2, "E_shnum", "Number of section header table entries",
	},
	{
		0x3E, 2, "E_shstrndx", "Index of section header table entry that holds the section names",
	},
}

var P_type_map = map[uint32]string{
	0:          "PT_NULL",
	1:          "PT_LOAD",
	2:          "PT_DYNAMIC",
	3:          "PT_INTERP",
	4:          "PT_NOTE",
	5:          "PT_SHLIB",
	6:          "PT_PHDR",
	7:          "PT_TLS",
	0x60000000: "PT_LOOS",
	0x6fffffff: "PT_HIOS",
	0x7fffffff: "PT_HIPROC",
	0x6474e550: "PT_GNU_EH_FRAME/PT_SUNW_EH_FRAME",
	0x6464e550: "PT_SUNW_UNWIND",
	0x6474e551: "PT_GNU_STACK",
	0x6474e552: "PT_GNU_RELRO",
	0x65a3dbe6: "PT_OPENBSD_RANDOMIZE",
	0x65a3dbe7: "PT_OPENBSD_WXNEEDED",
	0x65a41be6: "PT_OPENBSD_BOOTDATA",
	0x70000000: "PT_ARM_ARCHEXT/PT_MIPS_REGINFO,PT_LOPROC",
	0x70000001: "PT_ARM_EXIDX/PT_ARM_UNWIND/PT_MIPS_RTPROC",
	0x70000002: "PT_MIPS_OPTIONS",
	0x70000003: "PT_MIPS_ABIFLAGS",
};

var ProgramHeaderDescription = []HeaderDescriptionEntry{
	{0, 4, "P_type", "Segment identifier"},
	{4, 4, "P_flags", "Segement dependent flags"},
	{8, 8, "P_offset", "Segment offset in file image"},
	{0x10, 0x8, "P_vaddr", "Virtual address in memory"},
	{0x18, 0x8, "P_addr", "Physical address"},
	{0x20, 0x8, "P_filesz", "Size in file image"},
	{0x28, 0x8, "P_memsz", "Size in memory"},
	{0x30, 0x8, "P_align", "Alignment - 0 or 1 is no alignment. P_vaddr == P_offset % P_align"},
}

var SH_type_map = map[uint32]string{
	0:          "SHT_NULL",
	1:          "SHT_PROGBITS",
	2:          "SHT_SYMTAB",
	3:          "SHT_STRTAB",
	4:          "SHT_RELA",
	5:          "SHT_HASH",
	6:          "SHT_DYNAMIC",
	7:          "SHT_NOTE",
	8:          "SHT_NOBITS",
	9:          "SHT_REL",
	10:         "SHT_SHLIB",
	11:         "SHT_DYNSYM",
	14:         "SHT_INIT_ARRAY",
	15:         "SHT_FINI_ARRAY",
	16:         "SHT_PREINIT_ARRAY",
	17:         "SHT_GROUP",
	18:         "SHT_SYMTAB_SHNDX",
	0x60000000: "SHT_LOOS",
	0x60000001: "SHT_ANDROID_REL",
	0x60000002: "SHT_ANDROID_RELA",
	0x6fff4c00: "SHT_LLVM_ODRTAB",
	0x6fff4c01: "SHT_LLVM_LINKER_OPTIONS",
	0x6ffffff5: "SHT_GNU_ATTRIBUTES",
	0x6ffffff6: "SHT_GNU_HASH",
	0x6ffffffd: "SHT_GNU_verdef",
	0x6ffffffe: "SHT_GNU_verneed",
	0x6fffffff: "SHT_GNU_versym/SHT_HIOS",
	0x70000000: "SHT_LOPROC/SHT_HEX_ORDERED",
	0x70000001: "SHT_ARM_EXIDX/SHT_X86_64_UNWIND",
	0x70000002: "SHT_ARM_PREEMPTMAP",
	0x70000003: "SHT_ARM_ATTRIBUTES",
	0x70000004: "SHT_ARM_DEBUGOVERLAY",
	0x70000005: "SHT_ARM_OVERLAYSECTION",
	0x70000006: "SHT_MIPS_REGINFO",
	0x7000000d: "SHT_MIPS_OPTIONS",
	0x7000001e: "SHT_MIPS_DWARF",
	0x7000002a: "SHT_MIPS_ABIFLAGS",
	0x7fffffff: "SHT_HIPROC",
	0x80000000: "SHT_LOUSER",
	0xffffffff: "SHT_HIUSER",
};

var SH_name_map = map[string]string{
	".shstrtab" : "table with section names",
	".interp" : "name of the interpreter",
	".note.ABI-tag" : "ABI details",
	".note.gnu.build-id" : "Build id",
	".gnu.hash" : "symbol hash table",
	".dynsym" : "dynamic linking symbol table",
	".dynstr" : "strings needed for dynamic linking",
	".gnu.version" : "Symbol Version table",
	".gnu.version_r" : "Version requirements",
	".rela.dyn" : "relocation info for dynamically linked objects",
	".init" : "instruction for process initialization",
	".text" : "Executable code",
	".fini" : "instructions for process termination",
	".rodata" : "Read Only data",
	".eh_frame_hdr" : "pointer to .eh_frame (or lookup table)",
	".eh_frame" : "info for frame unwinding",
	".init_array" : "function pointers used during process start",
	".fini_array" : "function pointers used during process termination",
	".data.rel.ro" : "read only after relocation",
	".rela.plt" : "procedure linkage table for element to be relocated",
	".dynamic" : "dynamic linking info",
	".got" : "global offset table",
	".plt": "procedure linkage table",
	".data" : "Initialized data",
	".bss" : "uninitialized data",
	".comment" : "Version control information",
	".debug": "debugging symbols",
}
var SH_flags_map = map[uint64]string{
	0x1: "SHF_WRITE",

	0x2: "SHF_ALLOC",

	0x4: "SHF_EXECINSTR",

	0x10: "SHF_MERGE",

	0x20: "SHF_STRINGS",

	0x40: "SHF_INFO_LINK",

	0x80: "SHF_LINK_ORDER",

	0x100: "SHF_OS_NONCONFORMING",

	0x200: "SHF_GROUP",

	0x400: "SHF_TLS",

	0x800: "SHF_COMPRESSED",

	0x80000000: "SHF_EXCLUDE/SHF_MIPS_STRING",

	0x0ff00000: "SHF_MASKOS",

	0xf0000000: "SHF_MASKPROC",

	0x10000000: "XCORE_SHF_DP_SECTION/SHF_X86_64_LARGE/SHF_HEX_GPREL/SHF_MIPS_GPREL",

	0x20000000: "XCORE_SHF_CP_SECTION/SHF_MIPS_NAMES/SHF_MIPS_MERGE",

	0x01000000: "SHF_MIPS_NODUPES",

	0x04000000: "SHF_MIPS_LOCAL",

	0x08000000: "SHF_MIPS_NOSTRIP",

	0x40000000: "SHF_MIPS_ADDR",

	0x2000000: "SHF_ARM_PURECODE",
};

var SectionHeaderDescription = []HeaderDescriptionEntry{
	{0, 4, "SH_name", "Offset in .shstrtab to section name"},
	{0x4, 0x4, "SH_type", "Type identifier"},
	{0x8, 0x8, "SH_flags", "Section attributes"},
	{0x10, 0x8, "SH_addr", "Virtual address in memory if loaded"},
	{0x18, 0x8, "SH_offset", "Offset in file image"},
	{0x20, 0x8, "SH_size", "Size in file - may be 0"},
	{0x28, 0x4, "SH_link", "Section index"},
	{0x2C, 0x4, "SH_info", "Extra information"},
	{0x30, 0x8, "SH_addralign", "Alignment of section"},
	{0x38, 0x8, "SH_entsize", "Size for entries, if applicable, 0 otherwise"},
}

type ELFHeaderData struct {
	EI_MAG        uint32
	EI_CLASS      uint8
	EI_DATA       uint8
	EI_VERSION    uint8
	EI_OSABI      uint8
	EI_ABIVERSION uint8
	EI_PAD        uint64
	E_type        uint16
	E_machine     uint16
	E_version     uint32
	E_entry       uint64
	E_phoff       uint64
	E_shoff       uint64
	E_flags       uint32
	E_ehsize      uint16
	E_phentsize   uint16
	E_phnum       uint16
	E_shentsize   uint16
	E_shnum       uint16
	E_shstrndx    uint16
}

type ProgramHeaderData struct {
	P_type   uint32
	P_flags  uint32
	P_offset uint64
	P_vaddr  uint64
	P_addr   uint64
	P_filesz uint64
	P_memsz  uint64
	P_align  uint64
}

type SectionHeaderData struct {
	SH_name      uint32
	SH_type      uint32
	SH_flags     uint64
	SH_addr      uint64
	SH_offset    uint64
	SH_size      uint64
	SH_link      uint32
	SH_info      uint32
	SH_addralign uint64
	SH_entsize   uint64
}

func PrintHeaderDescription(headerDescription []HeaderDescriptionEntry, what string) {
	fmt.Printf("%s entries:\n", what)
	for _, header_entry := range headerDescription {
		fmt.Printf("Size: %d Offset: %x\tName: %s \t\t | %s\n",
			header_entry.size, header_entry.offset, header_entry.field, header_entry.description)
	}
}

func PrintProgramHeadersDescription() {
	return
}

func PrintSectionHeadersDescription() {
	return
}

func (eh *ELFHeaderData) String() string {
	r := strings.Builder{}
	r.WriteString(fmt.Sprintf("ELF Header:\n"))
	r.WriteString(fmt.Sprintf("EI_MAG: %#x\n", (*eh).EI_MAG))
	r.WriteString(fmt.Sprintf("EI_CLASS: %#x\n", (*eh).EI_CLASS))
	r.WriteString(fmt.Sprintf("EI_DATA: %#x\n", (*eh).EI_DATA))
	r.WriteString(fmt.Sprintf("EI_VERSION: %#x\n", (*eh).EI_VERSION))
	r.WriteString(fmt.Sprintf("EI_OSABI: %#x - %s\n", (*eh).EI_OSABI, E_OSABI_map[(*eh).EI_OSABI]))
	r.WriteString(fmt.Sprintf("EI_ABIVERSION: %#x\n", (*eh).EI_ABIVERSION))
	r.WriteString(fmt.Sprintf("EI_PAD: %#x\n", (*eh).EI_PAD))
	r.WriteString(fmt.Sprintf("E_type: %#x\n", (*eh).E_type))
	r.WriteString(fmt.Sprintf("E_machine: %#x - %s\n", (*eh).E_machine, E_machine_map[(*eh).E_machine]))
	r.WriteString(fmt.Sprintf("E_version: %#x\n", (*eh).E_version))
	r.WriteString(fmt.Sprintf("E_entry: %#x\n", (*eh).E_entry))
	r.WriteString(fmt.Sprintf("E_phoff: %#x\n", (*eh).E_phoff))
	r.WriteString(fmt.Sprintf("E_shoff: %#x\n", (*eh).E_shoff))
	r.WriteString(fmt.Sprintf("E_flags: %#x\n", (*eh).E_flags))
	r.WriteString(fmt.Sprintf("E_ehsize: %#x\n", (*eh).E_ehsize))
	r.WriteString(fmt.Sprintf("E_phentsize: %#x\n", (*eh).E_phentsize))
	r.WriteString(fmt.Sprintf("E_phnum: %#x\n", (*eh).E_phnum))
	r.WriteString(fmt.Sprintf("E_shentsize: %#x\n", (*eh).E_shentsize))
	r.WriteString(fmt.Sprintf("E_shnum: %#x\n", (*eh).E_shnum))
	r.WriteString(fmt.Sprintf("E_shstrndx: %#x\n", (*eh).E_shstrndx))
	return r.String()
}

func (eh *ProgramHeaderData) String() string {
	r := strings.Builder{}

	r.WriteString(fmt.Sprintf("Program Header:\n"))
	r.WriteString(fmt.Sprintf("P_type: %#x - %s\n", (*eh).P_type, P_type_map[(*eh).P_type]))
	r.WriteString(fmt.Sprintf("P_flags: %#x\n", (*eh).P_flags))
	r.WriteString(fmt.Sprintf("P_offset: %#x\n", (*eh).P_offset))
	r.WriteString(fmt.Sprintf("P_vaddr: %#x\n", (*eh).P_vaddr))
	r.WriteString(fmt.Sprintf("P_addr: %#x\n", (*eh).P_addr))
	r.WriteString(fmt.Sprintf("P_filesz: %#x\n", (*eh).P_filesz))
	r.WriteString(fmt.Sprintf("P_memsz: %#x\n", (*eh).P_memsz))
	r.WriteString(fmt.Sprintf("P_align: %#x\n", (*eh).P_align))
	return r.String()
}

func (eh *SectionHeaderData) String() string {
	r := strings.Builder{}

	r.WriteString(fmt.Sprintf("Section Header:\n"))
	r.WriteString(fmt.Sprintf("SH_name: %#x\n", (*eh).SH_name))
	r.WriteString(fmt.Sprintf("SH_type: %#x - %s\n", (*eh).SH_type, SH_type_map[(*eh).SH_type]))
	r.WriteString(fmt.Sprintf("SH_flags: %#x - %s\n", (*eh).SH_flags, SH_flags_map[(*eh).SH_flags]))
	r.WriteString(fmt.Sprintf("SH_addr: %#x\n", (*eh).SH_addr))
	r.WriteString(fmt.Sprintf("SH_offset: %#x\n", (*eh).SH_offset))
	r.WriteString(fmt.Sprintf("SH_size: %#x\n", (*eh).SH_size))
	r.WriteString(fmt.Sprintf("SH_link: %#x\n", (*eh).SH_link))
	r.WriteString(fmt.Sprintf("SH_info: %#x\n", (*eh).SH_info))
	r.WriteString(fmt.Sprintf("SH_addralign: %#x\n", (*eh).SH_addralign))
	r.WriteString(fmt.Sprintf("SH_entsize: %#x\n", (*eh).SH_entsize))
	return r.String()
}
