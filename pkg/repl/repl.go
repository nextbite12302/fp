package repl

import (
	"context"
	"errors"
	"fmt"
	"fp/pkg/fp"
	"sort"
)

type REPL interface {
	ReplyInput(ctx context.Context, input string) (output string, executed bool)
	ClearBuffer() (output string)
}

type fpRepl struct {
	runtime *fp.Runtime
	parser  *fp.Parser
	buffer  string
}

func (r *fpRepl) ReplyInput(ctx context.Context, input string) (output string, executed bool) {
	tokenList := fp.Tokenize(input)
	executed = false
	if len(tokenList) == 0 {
		executed = true
	} else {
		for _, token := range tokenList {
			expr := r.parser.Input(token)
			if expr != nil {
				executed = true

				lastFrame := make(fp.Frame).Update(r.runtime.Stack[len(r.runtime.Stack)-1])
				stackSize := len(r.runtime.Stack)
				output, err := r.runtime.Step(ctx, expr)
				if err != nil {
					if errors.Is(err, fp.InterruptError) {
						// reset stack size
						r.runtime.Stack = r.runtime.Stack[:stackSize-1]
						r.runtime.Stack = append(r.runtime.Stack, lastFrame)
						r.writeln("interrupted - stack was recovered")
					}
					r.writeln(err.Error())
					continue
				}
				r.write("%v\n", output)
			}
		}
	}
	return r.flush(), executed
}

func (r *fpRepl) ClearBuffer() (output string) {
	r.parser.Clear()
	r.writeln("(Control + C) to clear parser buffer, (Control + D) to exit")
	return r.flush()
}

func (r *fpRepl) flush() (output string) {
	output, r.buffer = r.buffer, ""
	return output
}

func (r *fpRepl) write(format string, a ...interface{}) {
	r.buffer += fmt.Sprintf(format, a...)
}
func (r *fpRepl) writeln(format string, a ...interface{}) {
	r.write(format+"\n", a...)
}

func NewFP(runtime *fp.Runtime) (repl REPL, welcome string) {
	r := &fpRepl{
		runtime: runtime,
		parser:  &fp.Parser{},
		buffer:  "",
	}
	r.writeln("welcome to fp repl! type function or module name for help")
	r.write("loaded modules: ")
	var funcNameList []string
	for k := range r.runtime.Stack[0] {
		funcNameList = append(funcNameList, string(k))
	}
	sort.Strings(funcNameList)
	for _, name := range funcNameList {
		r.write("%s ", name)
	}
	r.writeln("")
	return r, r.flush()
}
