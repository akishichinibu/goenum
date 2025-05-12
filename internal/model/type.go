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

type ArrayType struct {
	ElemType Type
}

var _ Type = &ArrayType{}

func NewTypeArray(elemType Type) Type {
	return &ArrayType{
		ElemType: elemType,
	}
}
