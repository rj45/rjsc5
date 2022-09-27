#!/bin/env ruby

projpath = File.absolute_path(File.join(__dir__, ".."))

testbench = File.read("#{projpath}/dig/testbench.dig")
testsglob = "#{projpath}/tests/*.hex"

Dir[testsglob].each do |filename|
  testname = File.basename(filename, ".hex")
  # puts testname

  testprog = "#{projpath}/scripts/digtestgen.rb"
  replacementtext = `ruby #{testprog} #{filename}`


  i = testbench.index("<string>#{testname}</string>")
  if i.nil?
    warn "Could not find #{testname}!!!"
    # print testbench
    # exit 1

  else
    ds = "<dataString>"

    first = testbench.index(ds, i)+ds.length

    last = testbench.index("</dataString>", i)

    texttoreplace = testbench[first,last-first]

    encodedreplacement = replacementtext.
      encode(:xml => :text).
      gsub("\"", "&quot;").
      gsub("\'", "&apos;")

    testbench = testbench.sub(texttoreplace, encodedreplacement)
  end
end

print testbench
