// Copyright (c) 2024 Ryan "rj45" Sanche
// MIT Licensed, see LICENSE file for details.

// This program generates a ROM file translating all 16-bit RISC-V compressed instructions
// on the address pins to 32-bit RISC-V instructions on the data pins. Essentially this
// decompresses the compressed instructions into plain RV32I instructions.
//
// All reserved, float and RV64+ instructions generate 0 instructions, expected to be
// interpreted as illegal instructions by the processor.
//
// This program was mostly written by AI, with some minor modifications by me. As such
// the instruction formats (especially the immediate bits) may not be accurate.
//
// TODO: Verify the correctness of the instruction formats & remove this disclaimer.
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ROMSize          = 1 << 16 // 2^16 = 65536 entries
	stackPointerReg  = 2       // x2 is the stack pointer
	returnAddressReg = 1       // x1 is the return address register
	zeroReg          = 0       // x0 is the zero register
)

// Instruction represents a 32-bit RISC-V instruction. Compressed instructions
// are decompressed into this format to make it easier to generate the 32-bit
// version.
type Instruction struct {
	Opcode uint32
	Rd     uint32
	Funct3 uint32
	Rs1    uint32
	Rs2    uint32
	Funct7 uint32
	Imm    int32

	// Compressed instruction fields
	Funct2 uint32
	Funct4 uint32
	Funct6 uint32
}

////////////////////////////////////////////////////////////////////////////////////////
// Instruction format helper structs and functions
////////////////////////////////////////////////////////////////////////////////////////

// BitField represents a field within an instruction
type BitField struct {
	Name     string // name of the field (for documentation)
	StartBit int    // LSB position
	Length   int    // length in bits
}

// BitFields represents all fields in an instruction format
type BitFields struct {
	Name   string     // name of the instruction format
	Fields []BitField // fields in order
}

// Extract a field from an instruction
func extractField(inst uint16, field BitField) uint32 {
	mask := uint16((1 << field.Length) - 1)
	return uint32((inst >> field.StartBit) & mask)
}

// ImmediateFormat defines how to assemble immediate values from multiple fields
type ImmediateFormat struct {
	Name   string
	Fields []struct {
		Source BitField // where to get bits from
		Target int      // where to put them in result
	}
	SignBit int // position of sign bit (-1 for unsigned)
	Width   int // total width of immediate including sign bit
}

// Assemble an immediate value from its parts
func assembleImmediate(compressed uint16, format ImmediateFormat) int32 {
	var result int32

	// Assemble the immediate from its parts
	for _, field := range format.Fields {
		value := extractField(compressed, field.Source)
		result |= int32(value) << field.Target
	}

	// Sign extend if this is a signed immediate
	if format.SignBit >= 0 {
		// Check if sign bit is set
		if result&(1<<format.SignBit) != 0 {
			// Create mask for upper bits that need to be set
			upperMask := ^((1 << format.Width) - 1)
			result |= int32(upperMask)
		}
	}

	return result
}

// Helper function to assemble a 32-bit instruction from fields
func assembleInstruction(fields BitFields, values map[string]uint32) uint32 {
	var result uint32
	for _, field := range fields.Fields {
		value := values[field.Name]
		mask := uint32((1 << field.Length) - 1)
		result |= (value & mask) << field.StartBit
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////////////
// Encode functions for different 32-bit RV32I instruction formats
////////////////////////////////////////////////////////////////////////////////////////

var (
	RFormat = BitFields{
		Name: "R",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"rd", 7, 5},
			{"funct3", 12, 3},
			{"rs1", 15, 5},
			{"rs2", 20, 5},
			{"funct7", 25, 7},
		},
	}

	IFormat = BitFields{
		Name: "I",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"rd", 7, 5},
			{"funct3", 12, 3},
			{"rs1", 15, 5},
			{"imm[11:0]", 20, 12},
		},
	}

	SFormat = BitFields{
		Name: "S",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"imm[4:0]", 7, 5},
			{"funct3", 12, 3},
			{"rs1", 15, 5},
			{"rs2", 20, 5},
			{"imm[11:5]", 25, 7},
		},
	}

	BFormat = BitFields{
		Name: "B",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"imm[11]", 7, 1},
			{"imm[4:1]", 8, 4},
			{"funct3", 12, 3},
			{"rs1", 15, 5},
			{"rs2", 20, 5},
			{"imm[10:5]", 25, 6},
			{"imm[12]", 31, 1},
		},
	}

	UFormat = BitFields{
		Name: "U",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"rd", 7, 5},
			{"imm[31:12]", 12, 20},
		},
	}

	JFormat = BitFields{
		Name: "J",
		Fields: []BitField{
			{"opcode", 0, 7},
			{"rd", 7, 5},
			{"imm[19:12]", 12, 8},
			{"imm[11]", 20, 1},
			{"imm[10:1]", 21, 10},
			{"imm[20]", 31, 1},
		},
	}
)

func (i *Instruction) EncodeR() uint32 {
	return assembleInstruction(RFormat, map[string]uint32{
		"opcode": i.Opcode,
		"rd":     i.Rd,
		"funct3": i.Funct3,
		"rs1":    i.Rs1,
		"rs2":    i.Rs2,
		"funct7": i.Funct7,
	})
}

func (i *Instruction) EncodeI() uint32 {
	return assembleInstruction(IFormat, map[string]uint32{
		"opcode":    i.Opcode,
		"rd":        i.Rd,
		"funct3":    i.Funct3,
		"rs1":       i.Rs1,
		"imm[11:0]": uint32(i.Imm),
	})
}

func (i *Instruction) EncodeS() uint32 {
	return assembleInstruction(SFormat, map[string]uint32{
		"opcode":    i.Opcode,
		"imm[4:0]":  uint32(i.Imm) & 0x1F,
		"funct3":    i.Funct3,
		"rs1":       i.Rs1,
		"rs2":       i.Rs2,
		"imm[11:5]": uint32(i.Imm>>5) & 0x7F,
	})
}

func (i *Instruction) EncodeB() uint32 {
	imm := uint32(i.Imm)
	return assembleInstruction(BFormat, map[string]uint32{
		"opcode":    i.Opcode,
		"imm[11]":   (imm >> 11) & 0x1,
		"imm[4:1]":  (imm >> 1) & 0xF,
		"funct3":    i.Funct3,
		"rs1":       i.Rs1,
		"rs2":       i.Rs2,
		"imm[10:5]": (imm >> 5) & 0x3F,
		"imm[12]":   (imm >> 12) & 0x1,
	})
}

func (i *Instruction) EncodeU() uint32 {
	return assembleInstruction(UFormat, map[string]uint32{
		"opcode":     i.Opcode,
		"rd":         i.Rd,
		"imm[31:12]": uint32(i.Imm),
	})
}

func (i *Instruction) EncodeJ() uint32 {
	imm := uint32(i.Imm)
	return assembleInstruction(JFormat, map[string]uint32{
		"opcode":     i.Opcode,
		"rd":         i.Rd,
		"imm[19:12]": (imm >> 12) & 0xFF,
		"imm[11]":    (imm >> 11) & 0x1,
		"imm[10:1]":  (imm >> 1) & 0x3FF,
		"imm[20]":    (imm >> 20) & 0x1,
	})
}

////////////////////////////////////////////////////////////////////////////////////////
// Decode functions for compressed instruction formats
////////////////////////////////////////////////////////////////////////////////////////

// Compressed instruction formats
var (
	CRFormat = BitFields{
		Name: "CR",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rs2", 2, 5},
			{"rd/rs1", 7, 5},
			{"funct4", 12, 4},
		},
	}

	CIFormat = BitFields{
		Name: "CI",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"imm[5]", 12, 1},
			{"imm[4:0]", 2, 5},
			{"rd/rs1", 7, 5},
			{"funct3", 13, 3},
		},
	}

	CSSFormat = BitFields{
		Name: "CSS",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rs2", 2, 5},
			{"imm[5:0]", 7, 6},
			{"funct3", 13, 3},
		},
	}

	CIWFormat = BitFields{
		Name: "CIW",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rd'", 2, 3},
			{"imm[5:4]", 11, 2},
			{"imm[9:6]", 7, 4},
			{"imm[2]", 6, 1},
			{"imm[3]", 5, 1},
			{"funct3", 13, 3},
		},
	}

	CLFormat = BitFields{
		Name: "CL",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rd'", 2, 3},
			{"imm[5]", 12, 1},
			{"imm[4:3]", 10, 2},
			{"imm[2]", 6, 1},
			{"rs1'", 7, 3},
			{"funct3", 13, 3},
		},
	}

	CSFormat = BitFields{
		Name: "CS",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rs2'", 2, 3},
			{"imm[5]", 12, 1},
			{"imm[4:3]", 10, 2},
			{"imm[2]", 6, 1},
			{"rs1'", 7, 3},
			{"funct3", 13, 3},
		},
	}

	CAFormat = BitFields{
		Name: "CA",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rs2'", 2, 3},
			{"funct2", 5, 2},
			{"rs1'/rd'", 7, 3},
			{"funct6", 10, 6},
		},
	}

	CBFormat = BitFields{
		Name: "CB",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"offset[7:6]", 2, 2},
			{"offset[2:1]", 4, 2},
			{"offset[5]", 6, 1},
			{"rs1'", 7, 3},
			{"offset[4:3]", 10, 2},
			{"funct3", 13, 3},
		},
	}

	CJFormat = BitFields{
		Name: "CJ",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"jump_target[10]", 12, 1},
			{"jump_target[9:8]", 10, 2},
			{"jump_target[7]", 9, 1},
			{"jump_target[6]", 8, 1},
			{"jump_target[5]", 7, 1},
			{"jump_target[4]", 6, 1},
			{"jump_target[3:1]", 3, 3},
			{"jump_target[0]", 2, 1},
			{"funct3", 13, 3},
		},
	}

	CLWSPFormat = BitFields{
		Name: "CLWSP",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rd", 7, 5},
			{"imm[4:2]", 4, 3},
			{"imm[7:6]", 12, 2},
			{"imm[5]", 2, 1},
			{"funct3", 13, 3},
		},
	}

	CSWSPFormat = BitFields{
		Name: "CSWSP",
		Fields: []BitField{
			{"opcode", 0, 2},
			{"rs2", 2, 5},
			{"imm[5:2]", 7, 4},
			{"imm[7:6]", 11, 2},
			{"funct3", 13, 3},
		},
	}
)

// Immediate formats for compressed instructions
var (
	CIWImmediate = ImmediateFormat{
		Name: "CIW immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[5:4]", 11, 2}, 6},
			{BitField{"imm[9:6]", 7, 4}, 2},
			{BitField{"imm[2]", 6, 1}, 1},
			{BitField{"imm[3]", 5, 1}, 2},
		},
		SignBit: -1, // unsigned
		Width:   10,
	}

	CLImmediate = ImmediateFormat{
		Name: "CL immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[5]", 12, 1}, 5},
			{BitField{"imm[4:3]", 10, 2}, 3},
			{BitField{"imm[2]", 6, 1}, 2},
		},
		SignBit: -1, // unsigned
		Width:   7,
	}

	// ... (continuing immediate formats)
	CSImmediate = ImmediateFormat{
		Name: "CS immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[5]", 12, 1}, 5},
			{BitField{"imm[4:3]", 10, 2}, 3},
			{BitField{"imm[2]", 6, 1}, 2},
		},
		SignBit: -1, // unsigned
		Width:   6,
	}

	CIImmediate = ImmediateFormat{
		Name: "CI immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[5]", 12, 1}, 5},
			{BitField{"imm[4:0]", 2, 5}, 0},
		},
		SignBit: 5, // signed, bit 5 is sign
		Width:   6,
	}

	CJImmediate = ImmediateFormat{
		Name: "CJ immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"jump_target[11]", 12, 1}, 11},
			{BitField{"jump_target[10]", 8, 1}, 10},
			{BitField{"jump_target[9:8]", 9, 2}, 8},
			{BitField{"jump_target[7]", 6, 1}, 7},
			{BitField{"jump_target[6]", 7, 1}, 6},
			{BitField{"jump_target[5]", 2, 1}, 5},
			{BitField{"jump_target[4]", 11, 1}, 4},
			{BitField{"jump_target[3:1]", 3, 3}, 1},
			{BitField{"jump_target[0]", 2, 1}, 0},
		},
		SignBit: 11, // signed
		Width:   12,
	}

	CBImmediate = ImmediateFormat{
		Name: "CB immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"offset[8]", 12, 1}, 8},
			{BitField{"offset[7:6]", 5, 2}, 6},
			{BitField{"offset[5]", 2, 1}, 5},
			{BitField{"offset[4:3]", 10, 2}, 3},
			{BitField{"offset[2:1]", 3, 2}, 1},
		},
		SignBit: 8, // signed, bit 8 is sign
		Width:   9,
	}

	CLWSPImmediate = ImmediateFormat{
		Name: "CLWSP immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[7:6]", 12, 2}, 6},
			{BitField{"imm[5]", 2, 1}, 5},
			{BitField{"imm[4:2]", 4, 3}, 2},
			// Note: bits [1:0] are implicitly 0 for word alignment
		},
		SignBit: -1, // unsigned
		Width:   8,
	}

	CSWSPImmediate = ImmediateFormat{
		Name: "CSWSP immediate",
		Fields: []struct {
			Source BitField
			Target int
		}{
			{BitField{"imm[7:6]", 11, 2}, 6},
			{BitField{"imm[5:2]", 7, 4}, 2},
			// Note: bits [1:0] are implicitly 0 for word alignment
		},
		SignBit: -1, // unsigned
		Width:   8,
	}
)

func decodeCR(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CRFormat.Fields[2]), // rd/rs1 field
		Rs1:    extractField(compressed, CRFormat.Fields[2]), // rd/rs1 field
		Rs2:    extractField(compressed, CRFormat.Fields[1]), // rs2 field
		Funct4: extractField(compressed, CRFormat.Fields[3]), // funct4 field
	}
}

func decodeCI(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CIFormat.Fields[3]), // rd/rs1 field
		Rs1:    extractField(compressed, CIFormat.Fields[3]), // rd/rs1 field
		Funct3: extractField(compressed, CIFormat.Fields[4]), // funct3 field
		Imm:    assembleImmediate(compressed, CIImmediate),
	}
}

func decodeCIW(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CIWFormat.Fields[1]) + 8, // rd' field
		Rs1:    stackPointerReg,
		Funct3: extractField(compressed, CIWFormat.Fields[6]), // funct3 field
		Imm:    assembleImmediate(compressed, CIWImmediate),
	}
}

func decodeCL(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CLFormat.Fields[1]) + 8, // rd' field
		Rs1:    extractField(compressed, CLFormat.Fields[5]) + 8, // rs1' field
		Funct3: extractField(compressed, CLFormat.Fields[6]),     // funct3 field
		Imm:    assembleImmediate(compressed, CLImmediate),
	}
}

func decodeCS(compressed uint16) Instruction {
	return Instruction{
		Rs1:    extractField(compressed, CSFormat.Fields[5]) + 8, // rs1' field
		Rs2:    extractField(compressed, CSFormat.Fields[1]) + 8, // rs2' field
		Funct3: extractField(compressed, CSFormat.Fields[6]),     // funct3 field
		Imm:    assembleImmediate(compressed, CSImmediate),
	}
}

func decodeCA(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CAFormat.Fields[3]) + 8, // rd' field
		Rs1:    extractField(compressed, CAFormat.Fields[3]) + 8, // rs1' field
		Rs2:    extractField(compressed, CAFormat.Fields[1]) + 8, // rs2' field
		Funct6: extractField(compressed, CAFormat.Fields[4]),     // funct6 field
		Funct2: extractField(compressed, CAFormat.Fields[2]),     // funct2 field
	}
}

func decodeCB(compressed uint16) Instruction {
	return Instruction{
		Rs1:    extractField(compressed, CBFormat.Fields[4]) + 8, // rs1' field
		Funct3: extractField(compressed, CBFormat.Fields[6]),     // funct3 field
		Imm:    assembleImmediate(compressed, CBImmediate),
	}
}

func decodeCJ(compressed uint16) Instruction {
	return Instruction{
		Funct3: extractField(compressed, CJFormat.Fields[9]), // funct3 field
		Imm:    assembleImmediate(compressed, CJImmediate),
	}
}

func decodeCLWSP(compressed uint16) Instruction {
	return Instruction{
		Rd:     extractField(compressed, CLWSPFormat.Fields[1]),
		Rs1:    stackPointerReg,
		Funct3: extractField(compressed, CLWSPFormat.Fields[5]),
		Imm:    assembleImmediate(compressed, CLWSPImmediate),
	}
}

func decodeCSWSP(compressed uint16) Instruction {
	return Instruction{
		Rs2:    extractField(compressed, CSWSPFormat.Fields[1]),
		Rs1:    stackPointerReg,
		Funct3: extractField(compressed, CSWSPFormat.Fields[4]),
		Imm:    assembleImmediate(compressed, CSWSPImmediate),
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Decompression functions (translate compressed instructions to 32-bit instructions)
////////////////////////////////////////////////////////////////////////////////////////

// TODO: verify the correctness of these functions, especially the immediate values

func decompressQuadrant0(compressed uint16) uint32 {
	funct3 := extractField(compressed, BitField{"funct3", 13, 3})

	var instr Instruction

	if compressed == 0 {
		return 0 // C.ILLEGAL
	}

	switch funct3 {
	case 0x0: // C.ADDI4SPN
		instr = decodeCIW(compressed)
		if instr.Imm == 0 {
			return 0 // reserved
		}
		instr.Rd = stackPointerReg
		instr.Opcode = 0x13 // ADDI
		instr.Funct3 = 0x0  // ADDI
		return instr.EncodeI()
	case 0x1: // C.FLD (not supported in RV32IMAC)
		return 0
	case 0x2: // C.LW
		instr = decodeCL(compressed)
		instr.Opcode = 0x03 // LOAD
		instr.Funct3 = 0x2  // LW
		return instr.EncodeI()
	case 0x3: // C.LD (RV64) or C.FLW (RV32FC, not supported in RV32IMAC)
		return 0
	case 0x4: // Reserved
		return 0
	case 0x5: // C.FSD (not supported in RV32IMAC)
		return 0
	case 0x6: // C.SW
		instr = decodeCS(compressed)
		instr.Opcode = 0x23 // STORE
		instr.Funct3 = 0x2  // SW
		return instr.EncodeS()
	case 0x7: // C.SD (RV64) or C.FSW (RV32FC, not supported in RV32IMAC)
		return 0
	}

	return 0 // Invalid instruction
}

func decompressQuadrant1(compressed uint16) uint32 {
	funct3 := extractField(compressed, BitField{"funct3", 13, 3})

	var instr Instruction

	switch funct3 {
	case 0x0: // C.ADDI
		instr = decodeCI(compressed)
		instr.Opcode = 0x13 // ADDI
		instr.Funct3 = 0x0  // ADDI
		return instr.EncodeI()
	case 0x1: // C.JAL (RV32) (or C.ADDIW for RV64+)
		instr = decodeCJ(compressed)
		instr.Opcode = 0x6F // JAL
		instr.Rd = returnAddressReg
		return instr.EncodeJ()
	case 0x2: // C.LI
		instr = decodeCI(compressed)
		instr.Opcode = 0x13 // ADDI
		instr.Funct3 = 0x0  // ADDI
		instr.Rs1 = zeroReg
		return instr.EncodeI()
	case 0x3: // C.ADDI16SP or C.LUI
		instr = decodeCI(compressed)
		if instr.Rd == 0 {
			return 0 // reserved
		} else if instr.Rd == 2 { // C.ADDI16SP
			instr.Opcode = 0x13 // ADDI
			instr.Funct3 = 0x0  // ADDI
			instr.Rs1 = stackPointerReg
			instr.Imm <<= 4 // scale immediate by 16
			return instr.EncodeI()
		} else { // C.LUI
			if instr.Imm == 0 {
				return 0 // reserved
			}
			instr.Opcode = 0x37 // LUI
			instr.Imm <<= 12    // shift immediate
			return instr.EncodeU()
		}
	case 0x4: // C.SRLI, C.SRAI, C.ANDI, C.SUB, C.XOR, C.OR, C.AND
		instr = decodeCB(compressed)
		funct2 := extractField(compressed, BitField{"funct2", 10, 2})
		switch funct2 {
		case 0x0: // C.SRLI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x5  // SRLI
			instr.Funct7 = 0x00
			instr.Rd = instr.Rs1
			return instr.EncodeI()
		case 0x1: // C.SRAI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x5  // SRAI
			instr.Funct7 = 0x20
			instr.Rd = instr.Rs1
			return instr.EncodeI()
		case 0x2: // C.ANDI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x7  // ANDI
			instr.Rd = instr.Rs1
			return instr.EncodeI()
		case 0x3: // C.SUB, C.XOR, C.OR, C.AND
			instr = decodeCA(compressed)
			instr.Opcode = 0x33 // OP
			switch instr.Funct2 {
			case 0x0: // C.SUB
				instr.Funct3 = 0x0 // SUB
				instr.Funct7 = 0x20
			case 0x1: // C.XOR
				instr.Funct3 = 0x4 // XOR
			case 0x2: // C.OR
				instr.Funct3 = 0x6 // OR
			case 0x3: // C.AND
				instr.Funct3 = 0x7 // AND
			}
			return instr.EncodeR()
		}
	case 0x5: // C.J
		instr = decodeCJ(compressed)
		instr.Opcode = 0x6F // JAL
		instr.Rd = zeroReg
		return instr.EncodeJ()
	case 0x6: // C.BEQZ
		instr = decodeCB(compressed)
		instr.Opcode = 0x63 // BRANCH
		instr.Funct3 = 0x0  // BEQ
		instr.Rs2 = zeroReg
		return instr.EncodeB()
	case 0x7: // C.BNEZ
		instr = decodeCB(compressed)
		instr.Opcode = 0x63 // BRANCH
		instr.Funct3 = 0x1  // BNE
		instr.Rs2 = zeroReg
		return instr.EncodeB()
	}

	return 0 // Invalid instruction
}

func decompressQuadrant2(compressed uint16) uint32 {
	funct3 := extractField(compressed, BitField{"funct3", 13, 3})

	var instr Instruction

	switch funct3 {
	case 0x0: // C.SLLI
		instr = decodeCI(compressed)
		if instr.Rd == 0 {
			return 0 // reserved
		}
		instr.Opcode = 0x13 // OP-IMM
		instr.Funct3 = 0x1  // SLLI
		instr.Rs1 = instr.Rd
		return instr.EncodeI()
	case 0x1: // C.FLDSP (not supported in RV32IMAC)
		return 0
	case 0x2: // C.LWSP
		instr = decodeCLWSP(compressed)
		if instr.Rd == 0 {
			return 0 // reserved
		}
		instr.Opcode = 0x03 // LOAD
		instr.Funct3 = 0x2  // LW
		return instr.EncodeI()
	case 0x3: // C.LDSP (RV64) or C.FLWSP (RV32FC, not supported in RV32IMAC)
		return 0
	case 0x4: // C.JR, C.MV, C.EBREAK, C.JALR, C.ADD
		instr = decodeCR(compressed)
		if instr.Rs2 == 0 {
			if instr.Rs1 == 0 {
				// C.EBREAK
				return 0x00100073 // EBREAK
			} else {
				// C.JR
				instr.Opcode = 0x67 // JALR
				instr.Rd = zeroReg
				instr.Funct3 = 0x0
				instr.Imm = 0
				return instr.EncodeI()
			}
		} else {
			if instr.Rs1 == 0 {
				// C.MV
				instr.Opcode = 0x33 // OP
				instr.Funct3 = 0x0  // ADD
				instr.Funct7 = 0x00
				instr.Rs1 = zeroReg
				return instr.EncodeR()
			} else {
				// C.ADD
				instr.Opcode = 0x33 // OP
				instr.Funct3 = 0x0  // ADD
				instr.Funct7 = 0x00
				instr.Rd = instr.Rs1
				return instr.EncodeR()
			}
		}
	case 0x5: // C.FSDSP (not supported in RV32IMAC)
		return 0
	case 0x6: // C.SWSP
		instr = decodeCSWSP(compressed)
		instr.Opcode = 0x23 // STORE
		instr.Funct3 = 0x2  // SW
		return instr.EncodeS()
	case 0x7: // C.SDSP (RV64) or C.FSWSP (RV32FC, not supported in RV32IMAC)
		return 0
	}

	return 0 // Invalid instruction
}

func decompressInstruction(compressed uint16) uint32 {
	// Extract opcode
	opcode := compressed & 0x3

	switch opcode {
	case 0x0:
		return decompressQuadrant0(compressed)
	case 0x1:
		return decompressQuadrant1(compressed)
	case 0x2:
		return decompressQuadrant2(compressed)
	default:
		// 0x3 is not a valid compressed instruction
		return 0xFFFFFFFF // Return an invalid instruction
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Main function
////////////////////////////////////////////////////////////////////////////////////////

func main() {
	rom := make([]uint32, ROMSize)

	// Iterate through all possible 16-bit instructions
	for i := 0; i < ROMSize; i++ {
		compressed := uint16(i)
		decompressed := decompressInstruction(compressed)
		rom[i] = decompressed
	}

	// Write ROM to file
	file, err := os.Create("dig/experiments/riscv_decompress_rom.bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, instruction := range rom {
		err := binary.Write(file, binary.LittleEndian, instruction)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("ROM file generated successfully.")
}
