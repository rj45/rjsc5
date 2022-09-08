import cocotb
from cocotb.triggers import Timer
from cocotb.triggers import RisingEdge
from cocotb.clock import Clock

async def reset_dut(dut):
    """Reset the cpu."""
    if dut.clk.value.binstr != "1":
        await RisingEdge(dut.clk)
    dut.reset.value = 1
    await RisingEdge(dut.clk)
    dut.reset.value = 0
    await Timer(1, units="step")

@cocotb.test()
async def test_reset_goes_low(dut):
    """Test reset goes low."""

    c = Clock(dut.clk, 1, 'ns')
    await cocotb.start(c.start())
    await reset_dut(dut)

    assert dut.reset.value == 0, "reset is not 0"


@cocotb.test()
async def test_regfile_forwarding(dut):
    """Test register file forwarding."""

    c = Clock(dut.clk, 1, 'ns')
    await cocotb.start(c.start())
    await reset_dut(dut)

    # enable clocks for both write and read
    dut.rw_clken.value = 1
    dut.ex_clken.value = 1

    # read and write to the same register, x1
    dut.de_rs.value = 1
    dut.rw_rd.value = 1

    # read the lower half of the register
    dut.de_half.value = 0

    # first write the upper half of the register
    # as zero
    dut.rw_half.value = 1
    dut.rw_result.value = 0
    await RisingEdge(dut.clk)

    # write the lower half, this should forward
    # the value so it's availabe on the next clk
    dut.rw_half.value = 0
    dut.rw_result.value = 0x1234
    await RisingEdge(dut.clk)

    # wait 1 step for the register to emit
    await Timer(1, units="step")

    # dut._log.info(f"value:{dut.ex_src.value.binstr}")
    assert dut.ex_src.value == 0x1234

    # extra clock so the test result is shown in
    # the VCD output
    await RisingEdge(dut.clk)


@cocotb.test()
async def test_regfile_registers(dut):
    """Test register file can use all registers."""

    c = Clock(dut.clk, 1, 'ns')
    await cocotb.start(c.start())
    await reset_dut(dut)

    # enable clock for writing
    dut.rw_clken.value = 1
    dut.ex_clken.value = 0

    value = 0

    # cycle through all 32 registers
    for reg in range(32):
        # cycle through each half
        for half in [0, 1]:
            value += 1

            # write the value to the register half
            dut.rw_rd.value = reg
            dut.rw_half.value = half
            dut.rw_result.value = value
            await RisingEdge(dut.clk)

    # switch to reading
    dut.rw_clken.value = 0
    dut.ex_clken.value = 1

    value = 0

    # cycle through all 32 registers
    for reg in range(32):
        # cycle through each half
        for half in [0, 1]:
            value += 1

            # emitted value is the lower 20 bits
            # so 4 bits of the upper half is sent
            # when the lower half is read
            expected = value
            if half == 0:
                expected = ((value+1) & 0xf) << 16 | value

            dut.de_half.value = half
            dut.de_rs.value = reg

            await RisingEdge(dut.clk)
            # wait 1 step for the register to emit
            await Timer(1, units="step")

            if reg == 0:
                # zero register should stay zero
                assert dut.ex_src.value == 0
            else:
                assert dut.ex_src.value == expected