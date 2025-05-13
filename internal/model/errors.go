package model

import (
	"fmt"
	"go/ast"
)

type InvalidEnumError struct {
	unit   *GenUnit
	decl   *ast.GenDecl
	reason string
}

var _ error = &InvalidEnumError{}

func (e *InvalidEnumError) Error() string {
	pos := e.unit.Package.Fset.Position(e.decl.Pos())

	return fmt.Sprintf("invalid enum declaration: %s:%d: %s", pos.Filename, pos.Line, e.reason)
}

type InvalidVariantError struct {
	unit     *GenUnit
	decl     *ast.TypeSpec
	funcType *ast.FuncType
	name     string
	reason   string
}

var _ error = &InvalidVariantError{}

func (e *InvalidVariantError) Error() string {
	pos := e.unit.Package.Fset.Position(e.decl.Pos())
	funcPos := e.unit.Package.Fset.Position(e.funcType.Pos())

	return fmt.Sprintf("invalid variant declaration: %s:%d: %s, %s, %s", pos.Filename, pos.Line, e.name, e.reason, funcPos)
}
