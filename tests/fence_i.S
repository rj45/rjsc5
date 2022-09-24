# See LICENSE for license details.

#*****************************************************************************
# fence_i.S
#-----------------------------------------------------------------------------

# Test self-modifying code and the fence.i instruction.

li a3, 111
lh a0, insn
lh a1, insn+2

# test I$ hit
.align 6
sh a0, 2f, t0
sh a1, 2f+2, t0
fence.i

la a5, 2f
jalr t1, a5, 0
test_2: li gp, 2; nop; li x7, ((444) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne a3, x7, fail;

 # test prefetcher hit
li a4, 100
1: addi a4, a4, -1
bnez a4, 1b

sh a0, 3f, t0
sh a1, 3f+2, t0
fence.i

.align 6
la a5, 3f
jalr t1, a5, 0
test_3: li gp, 3; nop; li x7, ((777) & ((1 << (__riscv_xlen - 1) << 1) - 1)); bne a3, x7, fail;

bne x0, gp, pass; fail: fence; 1: beqz gp, 1b; sll gp, gp, 1; or gp, gp, 1; li a7, 93; addi a0, gp, 0; ecall; pass: fence; li gp, 1; li a7, 93; li a0, 0; ecall

unimp


insn:
  addi a3, a3, 333

2: addi a3, a3, 222
jalr a5, t1, 0

3: addi a3, a3, 555
jalr a5, t1, 0

