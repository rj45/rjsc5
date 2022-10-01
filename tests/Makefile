tests = \
	simple \
	add addi \
	beq bne blt bge bltu bgeu \
	sub \
	and andi \
	or ori \
	xor xori \
	slt slti sltiu sltu \
	sll slli \
	sra srai \
	srl srli \
	lui \
	auipc \
	jal jalr \
	lw sw \
	# lb lbu lh lhu \
	# sb sh \
	# fence_i \

tests_dump = $(addsuffix .dump, $(tests))
tests_elf = $(addsuffix .elf, $(tests))
tests_bin = $(addsuffix .bin, $(tests))
tests_hex = $(addsuffix .hex, $(tests))

XLEN ?= 32
RISCV_PREFIX ?= riscv$(XLEN)-unknown-elf-
RISCV_GCC ?= $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS ?= -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles -ffreestanding -march=rv$(XLEN)i
RISCV_OBJDUMP ?= $(RISCV_PREFIX)objdump --disassemble-all --disassemble-zeroes --section=.text --section=.text.startup --section=.text.init --section=.data -M numeric
RISCV_OBJCOPY ?= $(RISCV_PREFIX)objcopy -O binary

default: all

all: $(tests_elf) $(tests_dump) $(tests_bin) $(tests_hex)

junk += $(tests_elf) $(tests_dump) $(tests_bin) $(tests_hex)

%.dump: %.elf
	$(RISCV_OBJDUMP) $< > $@

%.bin: %.elf
	$(RISCV_OBJCOPY) $< $@

%.hex: %.bin
	hexdump -e '"%08x " "\n"' -v $< > $@

%.elf: %.S
	$(RISCV_GCC) $(RISCV_GCC_OPTS) -T./link.ld $< -o $@

clean:
	rm -rf $(junk)
