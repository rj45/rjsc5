# Sept 24, 2022

Added the tests from the [RISC-V tests repo](https://github.com/riscv-software-src/riscv-tests).

The preprocessor was run over the test files to evaluate the macros used in order to make the code more readable. Only the tests are included and the linker script puts them at address zero. This makes it easier to start running them in a partially built processor.

The load/store tests are not yet included as I am not sure how to modify them for a harvard architecture. The final CPU will likely have a bus arbiter to allow data access to program memory but that will not initially be implemented. Could potentially just store the required data in the proper place that the tests expect using store instructions.

There is another repo with more comprehensive tests in it. That may be looked at later.

For reference, the CPP command is:

```sh
cpp -I ../macros/scalar/ -I../../env/p add.S
```
