/* verilator lint_off DECLFILENAME */

// Dual Port RAM for the register file
module dpram(
  input clk,
  input rw_clken,
  input [4:0] rw_rd,
  input [15:0] rw_result,
	input [4:0] de_rs,
  output [15:0] ex_src
);

reg [15:0] mem [31:0];

// distributed RAM, or unregistered block RAM
assign ex_src = mem[de_rs];

always_ff @(posedge clk) begin
  if (rw_clken) begin
		mem[rw_rd] <= rw_result;
  end
end

endmodule


// Dual port RAM for half of half of the register file
module regfile_dpram(
  input clk,
  input rw_clken,
  input [4:0] rw_rd,
  input [15:0] rw_result,
	input [4:0] de_rs,
  output [15:0] ex_src
);

wire [15:0] ex_src_;

dpram regs(
  .*,
  .ex_src(ex_src_)
);

wire zero_reg = de_rs == 5'd0;
wire forward_write = de_rs == rw_rd && rw_clken;

assign ex_src = zero_reg ? 16'd0 :
  (forward_write ? rw_result : ex_src_);

endmodule


// Half of the register file (one read port), with
// a 20 bit output for use by the address adder.
module regfile_half(
  input clk,
  input rw_clken,
  input rw_half,
  input [4:0] rw_rd,
  input [15:0] rw_result,
  input de_half,
	input [4:0] de_rs,
  input ex_clken,
  output reg [19:0] ex_src
);

wire rw_lclken = ~rw_half & rw_clken;
wire rw_hclken = rw_half & rw_clken;

wire [15:0] de_lsrc;
wire [15:0] de_hsrc;

regfile_dpram lregs(
  .*,
  .rw_clken(rw_lclken),
  .ex_src(de_lsrc)
);

regfile_dpram hregs(
  .*,
  .rw_clken(rw_hclken),
  .ex_src(de_hsrc)
);

always_ff @(posedge clk) begin
  if (ex_clken) begin
    // this might not be a good idea when we get to
    // handling hazards, but the ability for memory
    // load/stores to be 2 cyles instead of 4 is
    // appealing.
    ex_src <= de_half ?
      {4'd0, de_hsrc} :
      {de_hsrc[3:0], de_lsrc};
  end
end

endmodule
