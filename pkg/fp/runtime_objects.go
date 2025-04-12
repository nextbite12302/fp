package fp

import (
	"context"
	"fmt"
)

// types - TODO implement custom data types like Int, List, Dict

// Object : union - TODO : introduce new data types
type Object interface {
	String() string
	MustTypeObject() // for type-safety every Object must implement this
}

func getType(o Object) String {
	switch o.(type) {
	case Int:
		return "Int"
	case String:
		return "String"
	case Lambda:
		return "Lambda"
	case Module:
		return "Module"
	case List:
		return "List"
	case Dict:
		return "Dict"
	case Wildcard:
		return "Wildcard"
	case Unwrap:
		// unfortunately, one cannot use (type *) to get the type if unwrap since it will try to replace unwrap the next argument
		return "Unwrap"
	default:
		return "unknown"
	}
}

type Int int

func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}

func (i Int) MustTypeObject() {}

type Dict map[Object]Object

func (d Dict) String() string {

	s := ""
	s += "{"
	for k, v := range d {
		s += fmt.Sprintf("%s -> %s,", k.String(), v.String())
	}
	s += "}"
	return s
}

func (d Dict) MustTypeObject() {}

type Unwrap struct{}

func (u Unwrap) String() string {
	return "*"
}

func (u Unwrap) MustTypeObject() {}

type Wildcard struct{}

func (w Wildcard) String() string {
	return "_"
}

func (w Wildcard) MustTypeObject() {}

type String string

func (s String) String() string {
	return string(s)
}

func (s String) MustTypeObject() {}

type Lambda struct {
	Params []String `json:"params,omitempty"`
	Impl   Expr     `json:"impl,omitempty"`
	Frame  Frame    `json:"frame,omitempty"`
}

func (l Lambda) String() string {
	s := "(lambda "
	for _, param := range l.Params {
		s += param.String() + " "
	}
	s += l.Impl.String()
	s += ")"
	return s
}

func (l Lambda) MustTypeObject() {}

type Module struct {
	Name String `json:"name,omitempty"`
	Exec func(ctx context.Context, r *Runtime, expr LambdaExpr) (Object, error)
	Man  string `json:"man,omitempty"`
}

func (m Module) String() string {
	return m.Man
}

func (m Module) MustTypeObject() {}

type List []Object

func (l List) String() string {
	s := ""
	s += "["
	for _, obj := range l {
		s += fmt.Sprintf("%v,", obj)
	}
	s += "]"
	return s
}

func (l List) MustTypeObject() {}
