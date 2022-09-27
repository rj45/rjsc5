all:
	# todo

updatetests:
	ruby scripts/updatetests.rb > dig/testbench.dig.new
	mv dig/testbench.dig.new dig/testbench.dig
