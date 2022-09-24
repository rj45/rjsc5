# risc-v tests

From the [RISC-V tests repo](https://github.com/riscv-software-src/riscv-tests). See LICENSE for copyright info.

These are used to verify the processor works as intended.

For simplicity, they are run through `cpp` so they are macro free and easier to read.

They're also slightly modified to make it easier for bringing up a new processor by skipping the bring-up code. The ecall instruction can be used to determine if the test passed or not by looking at `a0`. `ecall` can be set up to act as a temporary halt instruction.

For implementation, work your way through `simple.S` then the tests for each instruction in `simple.S`. Then the rest of the instructions.
