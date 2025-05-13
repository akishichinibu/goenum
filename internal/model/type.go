package model

import (
	"go/ast"
)

type Type interface{}

type TypeDirect struct {
	Indent ast.Expr
}

var _ Type = &TypeDirect{}

func NewType(ident ast.Expr) (Type, error) {
	return &TypeDirect{
		Indent: ident,
	}, nil
}
