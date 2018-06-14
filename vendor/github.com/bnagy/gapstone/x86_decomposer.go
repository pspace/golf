/*
Gapstone is a Go binding for the Capstone disassembly library. For examples,
try reading the *_test.go files.

	Library Author: Nguyen Anh Quynh
	Binding Author: Ben Nagy
	License: BSD style - see LICENSE file for details
    (c) 2013 COSEINC. All Rights Reserved.
*/

package gapstone

// #cgo LDFLAGS: -lcapstone
// #cgo freebsd CFLAGS: -I/usr/local/include
// #cgo freebsd LDFLAGS: -L/usr/local/lib
// #include <stdlib.h>
// #include <capstone/capstone.h>
import "C"
import "unsafe"
import "reflect"

// import "fmt"

// Accessed via insn.X86.XXX
type X86Instruction struct {
	Prefix   []byte
	Opcode   []byte
	Rex      byte
	AddrSize byte
	ModRM    byte
	Sib      byte
	Disp     int32
	SibIndex uint
	SibScale int8
	SibBase  uint
	SseCC    uint
	AvxCC    uint
	AvxSAE   bool
	AvxRM    uint
	Operands []X86Operand
}

// Number of Operands of a given X86_OP_* type
func (insn X86Instruction) OpCount(optype uint) int {
	count := 0
	for _, op := range insn.Operands {
		if op.Type == optype {
			count++
		}
	}
	return count
}

type X86Operand struct {
	Type          uint // X86_OP_* - determines which field is set below
	Reg           uint
	Imm           int64
	FP            float64
	Mem           X86MemoryOperand
	Size          uint8
	AvxBcast      uint
	AvxZeroOpmask bool
}

type X86MemoryOperand struct {
	Segment uint
	Base    uint
	Index   uint
	Scale   int
	Disp    int64
}

func fillX86Header(raw C.cs_insn, insn *Instruction) {

	if raw.detail == nil {
		return
	}

	// Cast the cs_detail union
	cs_x86 := (*C.cs_x86)(unsafe.Pointer(&raw.detail.anon0[0]))

	// copy the prefix array to a new []byte
	pref := make([]byte, 4)
	for i := 0; i < 4; i++ {
		pref[i] = byte(cs_x86.prefix[i])
	}

	// Same for the opcode array
	opc := make([]byte, 4)
	for i := 0; i < 4; i++ {
		opc[i] = byte(cs_x86.opcode[i])
	}

	x86 := X86Instruction{
		Prefix:   pref,
		Opcode:   opc,
		Rex:      byte(cs_x86.rex),
		AddrSize: byte(cs_x86.addr_size),
		ModRM:    byte(cs_x86.modrm),
		Sib:      byte(cs_x86.sib),
		Disp:     int32(cs_x86.disp),
		SibIndex: uint(cs_x86.sib_index),
		SibScale: int8(cs_x86.sib_scale),
		SibBase:  uint(cs_x86.sib_base),
		SseCC:    uint(cs_x86.sse_cc),
		AvxCC:    uint(cs_x86.avx_cc),
		AvxSAE:   bool(cs_x86.avx_sae),
		AvxRM:    uint(cs_x86.avx_rm),
	}

	// Cast the op_info to a []C.cs_x86_op
	var ops []C.cs_x86_op
	oih := (*reflect.SliceHeader)(unsafe.Pointer(&ops))
	oih.Data = uintptr(unsafe.Pointer(&cs_x86.operands[0]))
	oih.Len = int(cs_x86.op_count)
	oih.Cap = int(cs_x86.op_count)

	// Create the Go object for each operand
	for _, cop := range ops {

		if cop._type == X86_OP_INVALID {
			break
		}

		gop := X86Operand{
			Type:          uint(cop._type),
			Size:          uint8(cop.size),
			AvxBcast:      uint(cop.avx_bcast),
			AvxZeroOpmask: bool(cop.avx_zero_opmask),
		}

		switch cop._type {
		// fake a union by setting only the correct struct member
		case X86_OP_IMM:
			gop.Imm = int64(*(*C.int64_t)(unsafe.Pointer(&cop.anon0[0])))
		case X86_OP_FP:
			gop.FP = float64(*(*C.double)(unsafe.Pointer(&cop.anon0[0])))
		case X86_OP_REG:
			gop.Reg = uint(*(*C.uint)(unsafe.Pointer(&cop.anon0[0])))
		case X86_OP_MEM:
			cmop := (*C.x86_op_mem)(unsafe.Pointer(&cop.anon0[0]))
			gop.Mem = X86MemoryOperand{
				Segment: uint(cmop.segment),
				Base:    uint(cmop.base),
				Index:   uint(cmop.index),
				Scale:   int(cmop.scale),
				Disp:    int64(cmop.disp),
			}
		}

		x86.Operands = append(x86.Operands, gop)
	}

	insn.X86 = &x86
}

func decomposeX86(raws []C.cs_insn) []Instruction {
	decomposed := []Instruction{}
	for _, raw := range raws {
		decomp := new(Instruction)
		fillGenericHeader(raw, decomp)
		fillX86Header(raw, decomp)
		decomposed = append(decomposed, *decomp)
	}
	return decomposed
}
