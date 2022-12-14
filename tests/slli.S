# See LICENSE for license details.

#*****************************************************************************
# slli.S
#-----------------------------------------------------------------------------

# Test slli instruction.

  #-------------------------------------------------------------
  # Arithmetic tests
  #-------------------------------------------------------------

  test_2: li gp, 2; li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((0) | (-(((0) >> 11) & 1) << 11));; li x7, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_3: li gp, 3; li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((1) | (-(((1) >> 11) & 1) << 11));; li x7, ((0x0000000000000002) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_4: li gp, 4; li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((7) | (-(((7) >> 11) & 1) << 11));; li x7, ((0x0000000000000080) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_5: li gp, 5; li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((14) | (-(((14) >> 11) & 1) << 11));; li x7, ((0x0000000000004000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_6: li gp, 6; li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((31) | (-(((31) >> 11) & 1) << 11));; li x7, ((0x0000000080000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  test_7: li gp, 7; li x1, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((0) | (-(((0) >> 11) & 1) << 11));; li x7, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_8: li gp, 8; li x1, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((1) | (-(((1) >> 11) & 1) << 11));; li x7, ((0xfffffffffffffffe) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_9: li gp, 9; li x1, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((7) | (-(((7) >> 11) & 1) << 11));; li x7, ((0xffffffffffffff80) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_10: li gp, 10; li x1, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((14) | (-(((14) >> 11) & 1) << 11));; li x7, ((0xffffffffffffc000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_11: li gp, 11; li x1, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((31) | (-(((31) >> 11) & 1) << 11));; li x7, ((0xffffffff80000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  test_12: li gp, 12; li x1, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((0) | (-(((0) >> 11) & 1) << 11));; li x7, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_13: li gp, 13; li x1, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((1) | (-(((1) >> 11) & 1) << 11));; li x7, ((0x0000000042424242) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_14: li gp, 14; li x1, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((7) | (-(((7) >> 11) & 1) << 11));; li x7, ((0x0000001090909080) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_15: li gp, 15; li x1, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((14) | (-(((14) >> 11) & 1) << 11));; li x7, ((0x0000084848484000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_16: li gp, 16; li x1, ((0x0000000021212121) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((31) | (-(((31) >> 11) & 1) << 11));; li x7, ((0x1090909080000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;







  #-------------------------------------------------------------
  # Source/Destination tests
  #-------------------------------------------------------------

  test_17: li gp, 17; li x1, ((0x00000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x1, x1, ((7) | (-(((7) >> 11) & 1) << 11));; li x7, ((0x00000080) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x1, x7, fail;;

  #-------------------------------------------------------------
  # Bypassing tests
  #-------------------------------------------------------------

  test_18: li gp, 18; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((7) | (-(((7) >> 11) & 1) << 11)); addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000000080) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;
  test_19: li gp, 19; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((14) | (-(((14) >> 11) & 1) << 11)); nop; addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000004000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;
  test_20: li gp, 20; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((31) | (-(((31) >> 11) & 1) << 11)); nop; nop; addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000080000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;

  test_21: li gp, 21; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x14, x1, ((7) | (-(((7) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000000080) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_22: li gp, 22; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); nop; slli x14, x1, ((14) | (-(((14) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000004000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_23: li gp, 23; li x4, 0; 1: li x1, ((0x0000000000000001) & ((1 << (__riscv_xlen - 1) << 1) - 1)); nop; nop; slli x14, x1, ((31) | (-(((31) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000080000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  test_24: li gp, 24; slli x1, x0, ((31) | (-(((31) >> 11) & 1) << 11));; li x7, ((0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x1, x7, fail;;
  test_25: li gp, 25; li x1, ((33) & ((1 << (__riscv_xlen - 1) << 1) - 1)); slli x0, x1, ((20) | (-(((20) >> 11) & 1) << 11));; li x7, ((0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x0, x7, fail;;

  bne x0, gp, pass; fail: fence; 1: beqz gp, 1b; sll gp, gp, 1; or gp, gp, 1; li a7, 93; addi a0, gp, 0; ecall; pass: fence; li gp, 1; li a7, 93; li a0, 0; ecall

unimp
