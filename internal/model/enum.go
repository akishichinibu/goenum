package model

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ValueType string

const ValueTypeString ValueType = "string"
const ValueTypeInt ValueType = "int"
const ValueTypeFloat ValueType = "float"

type Enum struct {
	Spec      *ast.TypeSpec
	Variants  []*Variant
	ValueType *ValueType
}

func (e *Enum) AddVariant(variant *Variant) error {
	for _, v := range e.Variants {
		if v.ValueType != variant.ValueType {
			return fmt.Errorf("variant %s has different value type from %s", variant.Name, v.Name)
		}
	}
	e.Variants = append(e.Variants, variant)
	return nil
}

type InvalidEnumError struct {
	decl   *ast.GenDecl
	reason string
}

var _ error = &InvalidEnumError{}

func (e *InvalidEnumError) Error() string {
	return fmt.Sprintf("invalid enum declaration: %+v, %s", e.decl, e.reason)
}

func parseDeclAsEnum(_ *ast.GenDecl, spec *ast.TypeSpec, iface *ast.InterfaceType) (e *Enum, err error) {
	e = &Enum{
		Spec:      spec,
		Variants:  make([]*Variant, 0),
		ValueType: nil,
	}

	for _, method := range iface.Methods.List {
		for _, name := range method.Names {
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				return nil, fmt.Errorf("invalid method type for %s: %T", name.Name, method.Type)
			}

			variants, err := newVariant(e, method, funcType)
			if err != nil {
				return nil, err
			}

			for _, variant := range variants {
				for _, field := range funcType.Params.List {
					for _, param := range field.Names {
						ft, err := NewType(field.Type)
						if err != nil {
							return nil, err
						}
						p := newParam(variant, param.Name, ft)
						variant.AddParam(p)
					}
				}
				if err := e.AddVariant(variant); err != nil {
					return nil, err
				}
			}
		}
	}

	return e, nil
}

func NewEnum(decl *ast.GenDecl) (ms []*Enum, err error) {
	if decl.Tok != token.TYPE {
		return nil, &InvalidEnumError{
			decl:   decl,
			reason: fmt.Sprintf("godantic enum should be declared with `%s`, got %s", token.TYPE, decl.Tok),
		}
	}

	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		iface, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		meta, err := parseDeclAsEnum(decl, typeSpec, iface)
		if err != nil {
			return nil, err
		}

		ms = append(ms, meta)
	}

	return ms, nil
}
