`include "cpu/rjsc5.sv"

module test_soc(
  input clk,
  input reset,

  // regfile_half
  input rw_clken,
  input rw_half,
  input [4:0] rw_rd,
  input [15:0] rw_result,
  input de_half,
	input [4:0] de_rs,
  input ex_clken,
  output reg [19:0] ex_src
);

`ifdef COCOTB_SIM
  initial begin
    $dumpfile ("rjsc5.vcd");
    $dumpvars (0, test_soc);
    #1;
  end
`endif

rjsc5 cpu(.*);

regfile_half rf(.*);

endmodule
