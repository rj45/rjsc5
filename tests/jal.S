# See LICENSE for license details.

#*****************************************************************************
# jal.S
#-----------------------------------------------------------------------------

# Test jal instruction.

  #-------------------------------------------------------------
  # Test 2: Basic test
  #-------------------------------------------------------------

test_2:
  li gp, 2
  li ra, 0

  jal x4, target_2
linkaddr_2:
  nop
  nop

  j fail

target_2:
  la x2, linkaddr_2
  bne x2, x4, fail

  #-------------------------------------------------------------
  # Test delay slot instructions not executed nor bypassed
  #-------------------------------------------------------------

  test_3: li gp, 3; li ra, 1; jal x0, 1f; addi ra, ra, 1; addi ra, ra, 1; addi ra, ra, 1; addi ra, ra, 1; 1: addi ra, ra, 1; addi ra, ra, 1;; li x7, ((3) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne ra, x7, fail;
# 50 "../rv64ui/jal.S"
  bne x0, gp, pass; fail: fence; 1: beqz gp, 1b; sll gp, gp, 1; or gp, gp, 1; li a7, 93; addi a0, gp, 0; ecall; pass: fence; li gp, 1; li a7, 93; li a0, 0; ecall

unimp

