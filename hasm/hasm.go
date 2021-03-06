package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func debug(x ...interface{}) {
	for i, v := range x {
		fmt.Printf("%v: %#v\n", i, v)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func assemble(st *ST, s string) (string, bool) {
	switch t := commandType(s); t {
	case A_COMMAND:
		sym := symbol(s)
		return "0" + st.value(sym), true
	case C_COMMAND:
		dmn := destMnemonic(s)
		cmn := compMnemonic(s)
		jmn := jumpMnemonic(s)
		return "111" + comp(cmn) + dest(dmn) + jump(jmn), true
	case L_COMMAND:
	case IGNORE:
	}
	return "", false
}

func buildTable(st *ST, s string, addr int) int {
	switch t := commandType(s); t {
	case L_COMMAND:
		st.addEntry(symbol(s), addr)
		return addr
	case A_COMMAND, C_COMMAND:
		return addr + 1
	case IGNORE:
	}
	return addr
}

func main() {
	fname := os.Args[1]
	fsplit := strings.Split(fname, ".")

	rfile, err := os.Open(fname)
	check(err)
	defer rfile.Close()

	wfile, err := os.Create(fsplit[0] + ".hack")
	check(err)
	defer wfile.Close()

	st := initST()
	addr := 0
	text := ""
	scanner := bufio.NewScanner(rfile)
        // First pass for building symbol table.
	for scanner.Scan() {
		text = remove(scanner.Text())
		addr = buildTable(&st, text, addr)
	}
	err = scanner.Err()
	check(err)

	rfile.Seek(0, 0)
	scanner = bufio.NewScanner(rfile)
	writer := bufio.NewWriter(wfile)
        // Second pass for generating binary.
	for scanner.Scan() {
		text = remove(scanner.Text())
		bin, isPrint := assemble(&st, text)
		if isPrint {
			fmt.Fprintln(writer, bin)
		}
	}
	err = scanner.Err()
	check(err)
	writer.Flush()
}
