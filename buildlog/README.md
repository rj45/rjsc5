# rjsc5 Build Log

## Sept 6, 2022 - Part 1

I have been designing a successor to [rj32](https://github.com/rj45/rj32) for a while now.

I chose RISC-V because what bogged me down about rj32 was the compiler side of things.

But I wanted to make it unique so it would be enough of a challenge to keep me interested. So, it's 16 bits internally.

A lot of thought has already gone into the design of the processor, and I started working on the verilog, but since I am new to verilog I made some mistakes with the initial prototype.

So here's a new start, and this time I will keep a build log so I can stop spamming the discord server I haunt with updates no one seems to care about :-)

## Sept 6, 2022 - Part 2

I chose cocotb to build testbenches with. It seems pretty simple and straightforward and building testbenches in Verilog is more than a little bit annoying. Maybe SystemVerilog has some good advances there. But cocotb is impressive so far, I like it.

## Sept 8, 2022

Added the [register file](../verilog/cpu/mod/regfile_half.sv) and wrestled with cocotb to get some unit tests around it.

[Read more](./2022-09-08-register-file.md)

## Sept 24, 2022

Added the tests from the [RISC-V tests repo](https://github.com/riscv-software-src/riscv-tests).

The preprocessor was run over the test files to evaluate the macros used in order to make the code more readable. Only the tests are included and the linker script puts them at address zero. This makes it easier to start running them in a partially built processor.

The load/store tests are not yet included as I am not sure how to modify them for a harvard architecture. The final CPU will likely have a bus arbiter to allow data access to program memory but that will not initially be implemented. Could potentially just store the required data in the proper place that the tests expect using store instructions.

There is another repo with more comprehensive tests in it. That may be looked at later.

For reference, the CPP command is:

```sh
cpp -I ../macros/scalar/ -I../../env/p add.S
```
