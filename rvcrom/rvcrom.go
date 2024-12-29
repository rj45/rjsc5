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
	// Additional fields for compressed instructions
	Funct4 uint32
	Funct6 uint32
}

////////////////////////////////////////////////////////////////////////////////////////
// Encode functions for different 32-bit RV32I instruction formats
////////////////////////////////////////////////////////////////////////////////////////

// TODO: verify the correctness of these functions

func (i *Instruction) EncodeR() uint32 {
	return i.Opcode | (i.Rd << 7) | (i.Funct3 << 12) | (i.Rs1 << 15) | (i.Rs2 << 20) | (i.Funct7 << 25)
}

func (i *Instruction) EncodeI() uint32 {
	return i.Opcode | (i.Rd << 7) | (i.Funct3 << 12) | (i.Rs1 << 15) | ((uint32(i.Imm) & 0xFFF) << 20)
}

func (i *Instruction) EncodeS() uint32 {
	imm := uint32(i.Imm)
	return i.Opcode | ((imm & 0x1F) << 7) | (i.Funct3 << 12) | (i.Rs1 << 15) | (i.Rs2 << 20) | ((imm & 0xFE0) << 20)
}

func (i *Instruction) EncodeB() uint32 {
	imm := uint32(i.Imm)
	return i.Opcode | ((imm & 0x800) >> 4) | ((imm & 0x1E) << 7) | (i.Funct3 << 12) | (i.Rs1 << 15) | (i.Rs2 << 20) | ((imm & 0x7E0) << 20) | ((imm & 0x1000) << 19)
}

func (i *Instruction) EncodeU() uint32 {
	return i.Opcode | (i.Rd << 7) | (uint32(i.Imm) & 0xFFFFF000)
}

func (i *Instruction) EncodeJ() uint32 {
	imm := uint32(i.Imm)
	return i.Opcode | (i.Rd << 7) | ((imm & 0xFF000) << 12) | ((imm & 0x800) << 9) | ((imm & 0x7FE) << 20) | ((imm & 0x100000) << 11)
}

////////////////////////////////////////////////////////////////////////////////////////
// Decode functions for compressed instruction formats
////////////////////////////////////////////////////////////////////////////////////////

// TODO: verify the correctness of these functions

func decodeCR(compressed uint16) Instruction {
	return Instruction{
		Rd:     uint32((compressed >> 7) & 0x1F),
		Rs1:    uint32((compressed >> 7) & 0x1F),
		Rs2:    uint32((compressed >> 2) & 0x1F),
		Funct4: uint32((compressed >> 12) & 0xF),
	}
}

func decodeCI(compressed uint16) Instruction {
	imm := int32(compressed>>2) & 0x1F
	imm |= (int32(compressed>>12) & 0x1) << 5
	if imm&0x20 != 0 {
		imm |= ^0x3F // Sign-extend
	}
	return Instruction{
		Rd:     uint32((compressed >> 7) & 0x1F),
		Rs1:    uint32((compressed >> 7) & 0x1F),
		Imm:    imm,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCSS(compressed uint16) Instruction {
	return Instruction{
		Rs1:    stackPointerReg,
		Rs2:    uint32((compressed >> 2) & 0x1F),
		Imm:    int32((compressed>>7)&0x3F) << 2,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCIW(compressed uint16) Instruction {
	return Instruction{
		Rd:     uint32((compressed>>2)&0x7) + 8,
		Rs1:    stackPointerReg,
		Imm:    int32(((compressed>>5)&0x1)|((compressed>>6)&0x1)<<1|((compressed>>7)&0xF)<<2|((compressed>>11)&0x3)<<6) << 2,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCL(compressed uint16) Instruction {
	return Instruction{
		Rd:     uint32((compressed>>2)&0x7) + 8,
		Rs1:    uint32((compressed>>7)&0x7) + 8,
		Imm:    int32(((compressed>>5)&0x1)|((compressed>>10)&0x7)<<1) << 2,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCS(compressed uint16) Instruction {
	return Instruction{
		Rs1:    uint32((compressed>>7)&0x7) + 8,
		Rs2:    uint32((compressed>>2)&0x7) + 8,
		Imm:    int32(((compressed>>5)&0x1)|((compressed>>10)&0x7)<<1) << 2,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCA(compressed uint16) Instruction {
	return Instruction{
		Rd:     uint32((compressed>>7)&0x7) + 8,
		Rs1:    uint32((compressed>>7)&0x7) + 8,
		Rs2:    uint32((compressed>>2)&0x7) + 8,
		Funct6: uint32((compressed >> 10) & 0x3F),
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCB(compressed uint16) Instruction {
	imm := int32(((compressed >> 2) & 0x1) | ((compressed>>3)&0x3)<<1 | ((compressed>>5)&0x3)<<3 | ((compressed>>10)&0x3)<<5)
	if imm&0x40 != 0 {
		imm |= ^0x7F // Sign-extend
	}
	return Instruction{
		Rs1:    uint32((compressed>>7)&0x7) + 8,
		Imm:    imm,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

func decodeCJ(compressed uint16) Instruction {
	imm := int32(((compressed >> 2) & 0x1) | ((compressed>>3)&0x7)<<1 | ((compressed>>6)&0x1)<<4 | ((compressed>>7)&0x1)<<5 | ((compressed>>8)&0x3)<<6 | ((compressed>>10)&0x1)<<8 | ((compressed>>11)&0x1)<<9 | ((compressed>>12)&0x1)<<10)
	if imm&0x400 != 0 {
		imm |= ^0x7FF // Sign-extend
	}
	return Instruction{
		Imm:    imm,
		Funct3: uint32((compressed >> 13) & 0x7),
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Decompression functions (translate compressed instructions to 32-bit instructions)
////////////////////////////////////////////////////////////////////////////////////////

// TODO: verify the correctness of these functions, especially the immediate values

func decompressQuadrant0(compressed uint16) uint32 {
	funct3 := (compressed >> 13) & 0x7

	var instr Instruction

	if compressed == 0 {
		return 0 // C.ILLEGAL
	}

	switch funct3 {
	case 0x0: // C.ADDI4SPN
		instr = decodeCIW(compressed)
		if instr.Imm == 0 {
			return 0 // reserved
		} else {
			instr.Opcode = 0x13 // ADDI
			instr.Rd = stackPointerReg
		}
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
	funct3 := (compressed >> 13) & 0x7

	var instr Instruction

	switch funct3 {
	case 0x0: // C.ADDI
		instr = decodeCI(compressed)
		if instr.Rd == 0 {
			return 0 // C.NOP (HINT)
		}
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
		funct2 := (compressed >> 10) & 0x3
		switch funct2 {
		case 0x0: // C.SRLI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x5  // SRLI
			instr.Funct7 = 0x00
			return instr.EncodeI()
		case 0x1: // C.SRAI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x5  // SRAI
			instr.Funct7 = 0x20
			return instr.EncodeI()
		case 0x2: // C.ANDI
			instr.Opcode = 0x13 // OP-IMM
			instr.Funct3 = 0x7  // ANDI
			return instr.EncodeI()
		case 0x3: // C.SUB, C.XOR, C.OR, C.AND
			instr = decodeCA(compressed)
			instr.Opcode = 0x33 // OP
			switch (compressed >> 5) & 0x3 {
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
	funct3 := (compressed >> 13) & 0x7

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
		instr = decodeCI(compressed)
		if instr.Rd == 0 {
			return 0 // reserved
		}
		instr.Opcode = 0x03 // LOAD
		instr.Funct3 = 0x2  // LW
		instr.Rs1 = stackPointerReg
		instr.Imm = (instr.Imm & 0x3) | ((instr.Imm & 0x1C) << 2) | ((instr.Imm & 0x20) << 4)
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
		instr = decodeCSS(compressed)
		instr.Opcode = 0x23 // STORE
		instr.Funct3 = 0x2  // SW
		instr.Rs1 = stackPointerReg
		instr.Imm = (instr.Imm & 0x3F) << 2
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
