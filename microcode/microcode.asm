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

; mux_lhs
LHS_PC = 0 << 12
LHS_RS1 = 1 << 12

; mux_rhs
RHS_RS2 = 0 << 13
RHS_IMM = 1 << 13

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
done EN_REG | LHS_RS1 | RHS_RS2 | FUNC_FUNCT3 | SUB_ALT | RES_ALU

shift:   ;  2: 00010
done EN_REG | LHS_RS1 | RHS_RS2 | FUNC_FUNCT3 | SUB_ALT | RES_ALU

mul:     ;  3: 00011
next 0 ; not implemented

div:     ;  4: 00100
next 0 ; not implemented

rem:     ;  5: 00101
next 0 ; not implemented

opi:     ;  6: 00110
done EN_REG | LHS_RS1 | RHS_IMM | FUNC_FUNCT3 | SUB_UCODE | IMM_I | RES_ALU

shifti:  ;  7: 00111
done EN_REG | LHS_RS1 | RHS_IMM | FUNC_FUNCT3 | SUB_UCODE | IMM_I | RES_ALU

compare: ;  8: 01000
next ACK_NEXT | LHS_RS1 | RHS_RS2 | FUNC_ADD | SUB_UCODE | COMPARE

branch:  ;  9: 01001
done LHS_PC | RHS_IMM | IMM_I | SB_RD | FUNC_ADD | SUB_UCODE | BR_BRANCH

jal:     ; 10: 01010
next 0 ; not implemented

jalr:    ; 11: 01011
next 0 ; not implemented

lui:     ; 12: 01100
done EN_REG | RES_IMM | IMM_U

auipc:   ; 13: 01101
next 0 ; not implemented

addr:    ; 14: 01110
next 0 ; not implemented

load:    ; 15: 01111
next 0 ; not implemented

store:   ; 16: 10000
next 0 ; not implemented

fence:   ; 17: 10001
done 0 ; nop

ecall:   ; 18: 10010
next 0 ; not implemented

ebreak:  ; 19: 10011
next 0 ; not implemented

mret:    ; 20: 10100
next 0 ; not implemented

wfi:     ; 21: 10101
next 0 ; not implemented

csrrw:   ; 22: 10110
next 0 ; not implemented

csrrs:   ; 23: 10111
next 0 ; not implemented

csrrc:   ; 24: 11000
next 0 ; not implemented