package avo

import (
	"fmt"
	"io"
	"strings"
)

// dot is the pesky unicode dot used in Go assembly.
const dot = "\u00b7"

type Printer interface {
	Print(*File) error
}

type GoPrinter struct {
	w         io.Writer
	by        string // generated by
	copyright []string
	err       error
}

func NewGoPrinter(w io.Writer) *GoPrinter {
	return &GoPrinter{
		w:  w,
		by: "avo",
	}
}

func (p *GoPrinter) SetGeneratedBy(by string) {
	p.by = by
}

func (p *GoPrinter) Print(f *File) error {
	p.header()

	for _, fn := range f.Functions {
		p.function(fn)
	}

	return p.err
}

func (p *GoPrinter) header() {
	p.generated()
	p.nl()
	p.incl("textflag.h")
	p.nl()
}

func (p *GoPrinter) generated() {
	p.comment(fmt.Sprintf("Code generated by %s. DO NOT EDIT.", p.by))
}

func (p *GoPrinter) incl(path string) {
	p.printf("#include \"%s\"\n", path)
}

func (p *GoPrinter) comment(line string) {
	p.multicomment([]string{line})
}

func (p *GoPrinter) multicomment(lines []string) {
	for _, line := range lines {
		p.printf("// %s\n", line)
	}
}

func (p *GoPrinter) function(f *Function) {
	p.printf("TEXT %s%s(SB),0,$%d-%d\n", dot, f.Name(), f.FrameBytes(), f.ArgumentBytes())

	for _, i := range f.inst {
		p.printf("\t%s\t%s\n", i.Opcode, joinOperands(i.Operands))
	}
}

func (p *GoPrinter) nl() {
	p.printf("\n")
}

func (p *GoPrinter) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(p.w, format, args...); err != nil {
		p.err = err
	}
}

func joinOperands(operands []Operand) string {
	asm := make([]string, len(operands))
	for i, op := range operands {
		asm[i] = op.Asm()
	}
	return strings.Join(asm, ", ")
}
