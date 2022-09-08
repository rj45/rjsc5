`include "mod/regfile_half.sv"

module rjsc5(
  input clk,
  input reset
);

wire rw_clken = 1'b0;
wire rw_half = 1'b0;
wire [4:0] rw_rd = 5'd0;
wire [15:0] rw_result = 16'd0;
wire de_half1 = 1'b0;
wire de_half2 = 1'b0;
wire [4:0] de_rs1 = 5'd0;
wire [4:0] de_rs2 = 5'd0;
wire ex_clken = 1'b0;

reg [19:0] ex_src1;
reg [19:0] ex_src2;

regfile_half regs1(
  .*,
  .de_rs(de_rs1),
  .de_half(de_half1),
  .ex_src(ex_src1)
);

regfile_half regs2(
  .*,
  .de_rs(de_rs2),
  .de_half(de_half2),
  .ex_src(ex_src2)
);


endmodule
