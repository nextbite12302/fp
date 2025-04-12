package fp

import (
	"encoding/json"
	"errors"
	"strconv"
)

// NewCoreRuntime - runtime + core control flow extensions
func NewCoreRuntime() *Runtime {
	return (&Runtime{
		parseLiteral: func(lit String) (Object, error) {
			if len(lit) == 0 {
				return nil, errors.New("empty literal")
			}
			if lit == "_" {
				return Wildcard{}, nil
			}
			if lit == "*" {
				return Unwrap{}, nil
			}
			if lit[0] == '"' && lit[len(lit)-1] == '"' {
				str := ""
				if err := json.Unmarshal([]byte(lit), &str); err != nil {
					return nil, err
				}
				return String(str), nil
			}
			i, err := strconv.Atoi(lit.String())
			return Int(i), err
		},
		Stack: []Frame{
			make(Frame),
		},
	}).
		LoadModule(letModule).
		LoadModule(delModule).
		LoadModule(lambdaModule).
		LoadModule(caseModule)
}

// NewBasicRuntime : NewCoreRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewCoreRuntime().
		LoadExtension(tailExtension).
		LoadExtension(addExtension).
		LoadExtension(subExtension).
		LoadExtension(signExtension)
}

// NewStdRuntime : NewCoreRuntime + standard functions
func NewStdRuntime() *Runtime {
	return NewBasicRuntime().
		LoadExtension(mulExtension).
		LoadExtension(divExtension).
		LoadExtension(modExtension).
		LoadExtension(printExtension).
		LoadExtension(listExtension).
		LoadExtension(appendExtension).
		LoadExtension(sliceExtension).
		LoadExtension(peekExtension).
		LoadExtension(lenExtension).
		LoadModule(mapModule).
		LoadExtension(typeExtension).
		LoadModule(stackModule).
		LoadModule(kaboomModule).
		LoadExtension(doomExtension).
		LoadExtension(timeExtension).
		LoadExtension(rangeExtension)
}
