#!/bin/env ruby

require "pathname"

# read the whole file specified in the first argument
hex = File.read(ARGV[0])

# split into lines
lines = hex.split("\n")

# split all the lines on spaces
machcode = lines.map {|x| x.split}

# flatten an array of arrays to a single large array
machcode = machcode.flatten

# add `0x` prefix to each hex number
machcode = machcode.map {|x| "0x#{x}"}

# print the template for digital test cases
puts <<~DONE
clk pass fail pc uop testcase ir

# #{ARGV[0]}
program(#{machcode.join(", ")})

let i = 0;
while(!(pass | fail | (i >= 2000)))
  let i = i + 1;
  0 0 0 x x x x
  1 x x x x x x
end while
0 1 0 x x x x
DONE
