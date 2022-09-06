# rjsc5 Build Log

## Sept 6, 2022 - Part 1

I have been designing a successor to [rj32](https://github.com/rj45/rj32) for a while now.

I chose RISC-V because what bogged me down about rj32 was the compiler side of things.

But I wanted to make it unique so it would be enough of a challenge to keep me interested. So, it's 16 bits internally.

A lot of thought has already gone into the design of the processor, and I started working on the verilog, but since I am new to verilog I made some mistakes with the initial prototype.

So here's a new start, and this time I will keep a build log so I can stop spamming the discord server I haunt with updates no one seems to care about :-)

## Sept 6, 2022 - Part 2

I chose cocotb to build testbenches with. It seems pretty simple and straightforward and building testbenches in Verilog is more than a little bit annoying. Maybe SystemVerilog has some good advances there. But cocotb is impressive so far, I like it.
