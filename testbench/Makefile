# Makefile

# defaults
SIM ?= icarus
TOPLEVEL_LANG ?= verilog

VERILOG_SOURCES += $(PWD)/../verilog/test_soc.sv

INCLUDES += -I$(PWD)/../verilog
INCLUDES += -I$(PWD)/../verilog/cpu
INCLUDES += -I$(PWD)/../verilog/cpu/mod

COMPILE_ARGS += $(INCLUDES)

TOPLEVEL = test_soc

MODULE = testbench

# include cocotb's make rules to take care of the simulator setup
include $(shell cocotb-config --makefiles)/Makefile.sim

.PHONEY: lint
lint:
	verilator --lint-only -Wall ${INCLUDES} ${VERILOG_SOURCES}
