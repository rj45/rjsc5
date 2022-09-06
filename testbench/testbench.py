import cocotb
from cocotb.triggers import Timer
from cocotb.triggers import RisingEdge
from cocotb.clock import Clock

async def reset_dut(dut):
    """Reset the cpu."""
    dut.reset.value = 1
    await RisingEdge(dut.clk)
    await RisingEdge(dut.clk)
    dut.reset.value = 0
    await Timer(1, units="step")

@cocotb.test(timeout_time=500, timeout_unit="ns")
async def test_reset_goes_low(dut):
    """Test reset goes low."""

    c = Clock(dut.clk, 1, 'ns')
    await cocotb.start(c.start())
    await reset_dut(dut)

    assert dut.reset.value == 0, "reset is not 0"