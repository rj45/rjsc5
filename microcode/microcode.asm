#bits 32


; ready for the next instruction or busy?
READY = 0 << 0
BUSY  = 1 << 0

; register enables
EN_PC = 1 << 6
EN_REG = 1 << 7
EN_STR = 1 << 8
EN_MAR = 1 << 9
EN_IR = 1 << 10

COMPARE = 1 << 11

; mux_addr
ADDR_PC = 0 << 12
ADDR_RS1 = 1 << 12

; mux_rhs
ALU_RS2 = 0 << 13
ALU_IMM = 1 << 13

; mux_func
FUNC_ADD = 0 << 14
FUNC_FUNCT3 = 1 << 14

; mux_ack
ACK_NEXT = 0 << 15
ACK_MEM = 1 << 15

; mux_sub
SUB_UCODE = 0 << 16
SUB_ALT = 1 << 16

; mux_sb
SB_RS2 = 0 << 17
SB_RD = 1 << 17

; mux_br
BR_NONE = 0 << 18
BR_JUMP = 1 << 18
BR_BRANCH = 2 << 18

; mux_imm
IMM_I = 0 << 20
IMM_U = 1 << 20
IMM_B = 2 << 20
IMM_J = 3 << 20

; mux_res
RES_PC = 0 << 22
RES_ALU = 1 << 22
RES_IMM = 2 << 22
RES_MEM = 3 << 22

#ruledef {
  next {value} => le((value | BUSY )`32)
  done {value} => le((value | READY | EN_PC | EN_IR )`32)
  custom {value} => le((value )`32)
}

#bankdef bank
{
    #addr 0
    #size 1<<6
    #outp 0
    #fill
}

boot:    ;  0: 00000 (boot and illegal instruction)
next EN_IR ; don't inc PC

op:      ;  1: 00001
done EN_REG | ALU_RS2 | FUNC_FUNCT3 | SUB_ALT | RES_ALU

slt:     ;  2: 00010
done EN_REG | ALU_RS2 | FUNC_FUNCT3 | SUB_UCODE | COMPARE | RES_ALU

shift:   ;  3: 00011
done EN_REG | ALU_RS2 | FUNC_FUNCT3 | SUB_ALT | RES_ALU

mul:     ;  4: 00100
next 0 ; not implemented

div:     ;  5: 00101
next 0 ; not implemented

rem:     ;  6: 00110
next 0 ; not implemented

opi:     ;  7: 00111
done EN_REG | ALU_IMM | FUNC_FUNCT3 | SUB_UCODE | IMM_I | RES_ALU

slti:    ;  8: 01000
done EN_REG | ALU_IMM | FUNC_FUNCT3 | SUB_UCODE | COMPARE | IMM_I | RES_ALU

shifti:  ;  9: 01001
done EN_REG | ALU_IMM | FUNC_FUNCT3 | SUB_UCODE | IMM_I | RES_ALU

branch:  ; 10: 01010
done ACK_NEXT | ADDR_PC | ALU_RS2 | IMM_B | SB_RD | FUNC_ADD | SUB_UCODE | COMPARE | BR_BRANCH

jal:     ; 11: 01011
next 0 ; not implemented

jalr:    ; 12: 01100
next 0 ; not implemented

lui:     ; 13: 01101
done EN_REG | RES_IMM | IMM_U

auipc:   ; 14: 01110
next 0 ; not implemented

addr:    ; 15: 01111
next 0 ; not implemented

load:    ; 16: 10000
next 0 ; not implemented

store:   ; 17: 10001
next 0 ; not implemented

fence:   ; 18: 10010
done 0 ; nop

ecall:   ; 19: 10011
next 0 ; not implemented

ebreak:  ; 20: 10100
next 0 ; not implemented

mret:    ; 21: 10101
next 0 ; not implemented

wfi:     ; 22: 10110
next 0 ; not implemented

csrrw:   ; 23: 10111
next 0 ; not implemented

csrrs:   ; 24: 11000
next 0 ; not implemented

csrrc:   ; 25: 11001
next 0 ; not implemented
