package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CodeWriter struct {
	input    *os.File
	output   *os.File
	basename string
	scanner  *bufio.Scanner
	writer   *bufio.Writer
}

func createCodeWriter(inname string) *CodeWriter {
	s := strings.Split(inname, ".")
	basename, extension := s[0], s[1]

	if extension != "vm" {
		panic("Input file should have .vm extension.")
	}

	input, err := os.Open(inname)
	check(err)

	outname := basename + ".asm"

	output, err := os.Create(outname)
	check(err)

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	return &CodeWriter{input, output, basename, scanner, writer}
}

func (cw *CodeWriter) close() {
	err := cw.scanner.Err()
	fmt.Println(cw)
	check(err)
	cw.writer.Flush()

	cw.input.Close()
	cw.output.Close()
}

func (cw *CodeWriter) text() string {
	return cw.scanner.Text()
}

func (cw *CodeWriter) scan() bool {
	return cw.scanner.Scan()
}

func writeArithmetic(arg1 ARG1) string {
	switch arg1 {
	case "add":
		return "M=D+M"
	case "sub":
		return "M=M-D"
	case "neg":
		return "M=-M"
	case "eq":
		return "JEQ"
	case "gt":
		return "JGT"
	case "lt":
		return "JLT"
	case "and":
		return "M=D&M"
	case "or":
		return "M=D|A"
	case "not":
		return "M=!M"
	default:
		return ""
	}
}

func writePushPop(c COMMAND_TYPE, seg string, idx int) string {
	// Stack pointer is hold at RAM[0].
	// Stack base starts from RAM[256].
	switch seg {
	case "push":
		return `@0
  M=D
  M=M+1
`
	case "pop":
		return `@0
  D=M
  M=M-1
  `
	default:
		return ""
	}
}
