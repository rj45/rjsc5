# See LICENSE for license details.

#*****************************************************************************
# xori.S
#-----------------------------------------------------------------------------

# Test xori instruction.

  #-------------------------------------------------------------
  # Logical tests
  #-------------------------------------------------------------

  test_2: li gp, 2; li x1, ((0x0000000000ff0f00) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0xf0f) | (-(((0xf0f) >> 11) & 1) << 11));; li x7, ((0xffffffffff00f00f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_3: li gp, 3; li x1, ((0x000000000ff00ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11));; li x7, ((0x000000000ff00f00) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_4: li gp, 4; li x1, ((0x0000000000ff08ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x70f) | (-(((0x70f) >> 11) & 1) << 11));; li x7, ((0x0000000000ff0ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_5: li gp, 5; li x1, ((0xfffffffff00ff00f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11));; li x7, ((0xfffffffff00ff0ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  #-------------------------------------------------------------
  # Source/Destination tests
  #-------------------------------------------------------------

  test_6: li gp, 6; li x1, ((0xffffffffff00f700) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x1, x1, ((0x70f) | (-(((0x70f) >> 11) & 1) << 11));; li x7, ((0xffffffffff00f00f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x1, x7, fail;;

   #-------------------------------------------------------------
  # Bypassing tests
  #-------------------------------------------------------------

  test_7: li gp, 7; li x4, 0; 1: li x1, ((0x000000000ff00ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11)); addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x000000000ff00f00) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;
  test_8: li gp, 8; li x4, 0; 1: li x1, ((0x0000000000ff08ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x70f) | (-(((0x70f) >> 11) & 1) << 11)); nop; addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000ff0ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;
  test_9: li gp, 9; li x4, 0; 1: li x1, ((0xfffffffff00ff00f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11)); nop; nop; addi x6, x14, 0; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0xfffffffff00ff0ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x6, x7, fail;;

  test_10: li gp, 10; li x4, 0; 1: li x1, ((0x000000000ff00ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x000000000ff00f00) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_11: li gp, 11; li x4, 0; 1: li x1, ((0x0000000000ff0fff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); nop; xori x14, x1, ((0x00f) | (-(((0x00f) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0x0000000000ff0ff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_12: li gp, 12; li x4, 0; 1: li x1, ((0xfffffffff00ff00f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); nop; nop; xori x14, x1, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11)); addi x4, x4, 1; li x5, 2; bne x4, x5, 1b; li x7, ((0xfffffffff00ff0ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  test_13: li gp, 13; xori x1, x0, ((0x0f0) | (-(((0x0f0) >> 11) & 1) << 11));; li x7, ((0x0f0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x1, x7, fail;;
  test_14: li gp, 14; li x1, ((0x00ff00ff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); xori x0, x1, ((0x70f) | (-(((0x70f) >> 11) & 1) << 11));; li x7, ((0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x0, x7, fail;;

  bne x0, gp, pass; fail: fence; 1: beqz gp, 1b; sll gp, gp, 1; or gp, gp, 1; li a7, 93; addi a0, gp, 0; ecall; pass: fence; li gp, 1; li a7, 93; li a0, 0; ecall

unimp

