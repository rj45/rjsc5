# See LICENSE for license details.

#*****************************************************************************
# lb.S
#-----------------------------------------------------------------------------

# for harvard, write some test data to memory
la x1, tdat
li x14, 0x0ff000ff
sw x14, 0(x1)

# Test lb instruction.

  #-------------------------------------------------------------
  # Basic tests
  #-------------------------------------------------------------

  test_2: li gp, 2; li x15, 0xffffffffffffffff; la x1, tdat; lb x14, 0(x1);; li x7, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_3: li gp, 3; li x15, 0x0000000000000000; la x1, tdat; lb x14, 1(x1);; li x7, ((0x0000000000000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_4: li gp, 4; li x15, 0xfffffffffffffff0; la x1, tdat; lb x14, 2(x1);; li x7, ((0xfffffffffffffff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_5: li gp, 5; li x15, 0x000000000000000f; la x1, tdat; lb x14, 3(x1);; li x7, ((0x000000000000000f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  # Test with negative offset

  test_6: li gp, 6; li x15, 0xffffffffffffffff; la x1, tdat4; lb x14, -3(x1);; li x7, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_7: li gp, 7; li x15, 0x0000000000000000; la x1, tdat4; lb x14, -2(x1);; li x7, ((0x0000000000000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_8: li gp, 8; li x15, 0xfffffffffffffff0; la x1, tdat4; lb x14, -1(x1);; li x7, ((0xfffffffffffffff0) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;
  test_9: li gp, 9; li x15, 0x000000000000000f; la x1, tdat4; lb x14, 0(x1);; li x7, ((0x000000000000000f) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x14, x7, fail;;

  # Test with a negative base

  test_10: li gp, 10; la x1, tdat; addi x1, x1, -32; lb x5, 32(x1);; li x7, ((0xffffffffffffffff) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x5, x7, fail;





  # Test with unaligned base

  test_11: li gp, 11; la x1, tdat; addi x1, x1, -6; lb x5, 7(x1);; li x7, ((0x0000000000000000) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x5, x7, fail;





  #-------------------------------------------------------------
  # Bypassing tests
  #-------------------------------------------------------------

  test_12: li gp, 12; li x4, 0; 1: la x1, tdat2; lb x14, 1(x1); addi x6, x14, 0; li x7, 0xfffffffffffffff0; bne x6, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;;
  test_13: li gp, 13; li x4, 0; 1: la x1, tdat3; lb x14, 1(x1); nop; addi x6, x14, 0; li x7, 0x000000000000000f; bne x6, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;;
  test_14: li gp, 14; li x4, 0; 1: la x1, tdat1; lb x14, 1(x1); nop; nop; addi x6, x14, 0; li x7, 0x0000000000000000; bne x6, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;;

  test_15: li gp, 15; li x4, 0; 1: la x1, tdat2; lb x14, 1(x1); li x7, 0xfffffffffffffff0; bne x14, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;
  test_16: li gp, 16; li x4, 0; 1: la x1, tdat3; nop; lb x14, 1(x1); li x7, 0x000000000000000f; bne x14, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;
  test_17: li gp, 17; li x4, 0; 1: la x1, tdat1; nop; nop; lb x14, 1(x1); li x7, 0x0000000000000000; bne x14, x7, fail; addi x4, x4, 1; li x5, 2; bne x4, x5, 1b;

  #-------------------------------------------------------------
  # Test write-after-write hazard
  #-------------------------------------------------------------

  test_18: li gp, 18; la x5, tdat; lb x2, 0(x5); li x2, 2;; li x7, ((2) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x2, x7, fail;





  test_19: li gp, 19; la x5, tdat; lb x2, 0(x5); nop; li x2, 2;; li x7, ((2) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne x2, x7, fail;






  bne x0, gp, pass; fail: fence; 1: beqz gp, 1b; sll gp, gp, 1; or gp, gp, 1; li a7, 93; addi a0, gp, 0; ecall; pass: fence; li gp, 1; li a7, 93; li a0, 0; ecall

unimp

  .data
tdat:
tdat1: .byte 0xff
tdat2: .byte 0x00
tdat3: .byte 0xf0
tdat4: .byte 0x0f

