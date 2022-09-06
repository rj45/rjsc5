# Makefile

# defaults
SIM ?= icarus
TOPLEVEL_LANG ?= verilog

VERILOG_SOURCES += $(PWD)/../verilog/test_soc.sv

COMPILE_ARGS += -I$(PWD)/../verilog
COMPILE_ARGS += -I$(PWD)/../verilog/cpu

TOPLEVEL = test_soc

MODULE = testbench

# include cocotb's make rules to take care of the simulator setup
include $(shell cocotb-config --makefiles)/Makefile.sim