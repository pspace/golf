package disassembler

import (
	"github.com/bnagy/gapstone"
	"fmt"
	"github.com/pspace/golf/helper"
)

type DisassemblerSpecs struct {
	Arch int
	Mode uint

}

func Disassemble(spec DisassemblerSpecs, code *[]byte) helper.ErrorCode{
	engine, err := gapstone.New(spec.Arch, spec.Mode)

	if err != nil {
		fmt.Errorf("Error while initializing Capstone Engine: %d", err)
		return helper.ECAPSTONE
	}

	instructions, err := engine.Disasm(*code, 0,0)

	if err != nil {
		fmt.Errorf("Error while disassembling instructions: %d", err)
		return helper.EDISASSEMBLE
	}

	defer engine.Close()

	for _, i := range instructions{
		fmt.Printf("%s %s\n", i.Mnemonic, i.OpStr)
	}

	return 0
}