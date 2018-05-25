package elf

import (
	"fmt"
)

var sizeToTypeMap = map[int]string{
	0: "uint8",
	1: "uint8",
	2: "uint16",
	4: "uint32",
	7: "uint64",
	8: "uint64",
}

var sizeToConverterMap = map[int]string{
	0: "",
	1: "",
	2: "ELFByteOrder.Uint16",
	4: "ELFByteOrder.Uint32",
	7: "ELFByteOrder.Uint64",
	8: "ELFByteOrder.Uint64",
}

func generateHeaderStruct(headerDescription []HeaderDescriptionEntry, what string) {
	fmt.Printf("type %sHeaderData struct {\n", what)
	for _, entry := range headerDescription {
		fmt.Printf("%s\t%s\n", entry.field, sizeToTypeMap[int(entry.size)])
	}
	fmt.Printf("}")
}

func generateHeaderParser(headerDescription []HeaderDescriptionEntry, what string) {
	fmt.Printf("func Parse%sHeader(input *os.File) *%sHeaderData {\n"+
		"headerData := %sHeaderData{}\nvar data []byte\n", what, what, what)

	for _, entry := range headerDescription {
		if entry.size == 0 {
			continue
		}
		fmt.Printf("data = helper.ReadNextBytesFromFile(input, %d)\n", entry.size)
		if 1 < entry.size {

			fmt.Printf("headerData.%s = %s(data)\n", entry.field, sizeToConverterMap[int(entry.size)])
		} else {
			fmt.Printf("headerData.%s = data[0]\n", entry.field)
			if "EI_DATA" == entry.field {
				fmt.Printf("setByteOrder(data[0])\n")
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("return &headerData\n}\n")
}

func generateHeaderStringer(headerDescription []HeaderDescriptionEntry, what string) {
	fmt.Printf("func (eh *%sHeaderData) String() string {\n"+
		"r := strings.Builder{}\n"+
		"r.WriteString(fmt.Sprintf(\"%s Header:\\n\"))\n", what, what)

	for _, entry := range headerDescription {

		fmt.Printf(" r.WriteString(fmt.Sprintf(\"%s: %%#x\\n\", (*eh).%s))", entry.field, entry.field)
		fmt.Printf("\n")
	}

	fmt.Printf("return r.String()\n}\n")
}

func generateHeaderComponents(headerDescription []HeaderDescriptionEntry, what string) {
	fmt.Printf("\n-------- Generating %s header code --------\n", what)
	generateHeaderStruct(headerDescription, what)
	fmt.Printf("\n")
	generateHeaderParser(headerDescription, what)
	fmt.Printf("\n")
	generateHeaderStringer(headerDescription, what)
}

func GenerateELFHeaderComponents() {
	generateHeaderComponents(ELFHeaderDescription, "ELF")
}

func GenerateSectionHeaderComponents() {
	generateHeaderComponents(SectionHeaderDescription, "Section")

}

func GenerateProgramHeaderComponents() {
	generateHeaderComponents(ProgramHeaderDescription, "Program")

}

func GenerateAll(){
	GenerateELFHeaderComponents()
	GenerateProgramHeaderComponents()
	GenerateSectionHeaderComponents()
}
