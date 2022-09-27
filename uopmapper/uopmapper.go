package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	f, _ := os.Open("microcode/opcodes.csv")
	r := csv.NewReader(f)

	of, _ := os.Create("microcode/opmatches.csv")

	records, _ := r.ReadAll()

	ops := map[string]int{
		"nop": 0,
	}
	opnames := []string{"nop"}
	num := 1

	for _, row := range records[1:] {
		if _, found := ops[row[0]]; !found {
			ops[row[0]] = num
			opnames = append(opnames, row[0])
			num++
		}
	}

	opbits := map[string]int{}
	opbitcols := []string{}

	for i, col := range records[0][1:] {
		opbits[col] = i
	}

	for bit, col := range opbits {
		found := false
		for _, row := range records[1:] {
			if row[col+1] != "" {
				found = true
				break
			}
		}
		if !found {
			delete(opbits, bit)
		} else {
			opbitcols = append(opbitcols, bit)
		}
	}
	sort.Slice(opbitcols, func(i, j int) bool {
		return opbits[opbitcols[i]] < opbits[opbitcols[j]]
	})

	inwidth := len(opbits)
	outwidth := 0
	for (1 << outwidth) < num {
		outwidth++
	}
	ignore := ""

	for _, name := range opbitcols {
		fmt.Fprintf(of, "i%s,", name)
	}
	for i := 0; i < outwidth; i++ {
		if i != 0 {
			ignore += ","
		}
		fmt.Fprintf(of, ",o%d", i)
		ignore += "0"
	}
	fmt.Fprintln(of)

	for i := 0; i < (1 << inwidth); i++ {
		pattern := fmt.Sprintf("%b", i)
		for len(pattern) < inwidth {
			pattern = "0" + pattern
		}

		patterncsv := strings.Join(strings.Split(pattern, ""), ",")

		found := false
		for _, row := range records[1:] {
			match := true
			for i, bit := range opbitcols {
				col := opbits[bit]
				val := row[col+1]
				if val != "" && val != string(pattern[i]) {
					match = false
				}
			}

			if match {
				uopcode := fmt.Sprintf("%b", ops[row[0]])
				for len(uopcode) < outwidth {
					uopcode = "0" + uopcode
				}

				fmt.Fprintf(of, "%s,,%s\n", patterncsv, strings.Join(strings.Split(uopcode, ""), ","))
				found = true
				break
			}
		}

		if !found {
			fmt.Fprintf(of, "%s,,%s\n", patterncsv, ignore)
		}

	}

	// fmt.Println(opbitcols)

	// fmt.Println(outwidth, inwidth, num)
	of.Close()

	for i, name := range opnames {
		fmt.Printf("%-8s ; %2d: %05b\n", name+":", i, i)
	}
}
