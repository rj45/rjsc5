`include "cpu/rjsc5.sv"

module test_soc(
  input clk,
  input reset
);

`ifdef COCOTB_SIM
  initial begin
    $dumpfile ("rjsc5.vcd");
    $dumpvars (0, test_soc);
    #1;
  end
`endif

rjsc5 cpu(.*);

endmodule