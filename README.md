# rjsc5

rjsc5 (reesk-five) is a 16-bit RISC-V CPU. Because 32 bits is already done many times.

## Goal

Build a 16-bit RISC-V CPU that's feasible to build out of discrete logic chips (74xxx series ICs), though this specific project is to build it in an FPGA with verilog.

## Build Log

You can read my [build log here](buildlog/README.md).

## Specs

* RV32I
* 32-bit registers
* 16-bit ALU and data paths
* Pipelined (4? stages)
* Memory
  * Harvard Architecture
  * 16 bit data bus on both I & D memory ports
  * Both memory ports:
    * Can stall the cpu (supporting slow memory/IO)
    * Can be tied to the same bus with a bus arbiter (Von Neumann Architecture)

### Stretch Goals

* Virtual Block Interface inspired virtual memory
* Unix-like OS built in Rust
* Retro graphics through a Video Display Processor
* Retro sounding audio
* M extension
* C extension
