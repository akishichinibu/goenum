package model

import (
	"fmt"
	"go/ast"
)

type Variant struct {
	Enum *Enum
	Name string

	Params    []*Param
	ValueType Type
}

func (v *Variant) AddParam(param *Param) error {
	v.Params = append(v.Params, param)
	return nil
}

func newVariant(enum *Enum, field *ast.Field, funcType *ast.FuncType) (v []*Variant, err error) {
	for _, name := range field.Names {
		name := name.Name

		if enum.ValueType != nil {
			if len(funcType.Params.List) != 0 {
				return nil, fmt.Errorf("strenum %s should not have parameters", name)
			}
			if funcType.Results == nil {
				return nil, fmt.Errorf("strenum %s should have a return value", name)
			}
			if len(funcType.Results.List) != 1 {
				return nil, fmt.Errorf("strenum %s should have only one return value", name)
			}
		} else {
			if funcType.Results != nil {
				return nil, fmt.Errorf("enum %s should not have a return value", name)
			}
		}

		variant := &Variant{
			Enum:   enum,
			Name:   name,
			Params: make([]*Param, 0),
		}

		v = append(v, variant)
	}
	return v, nil
}
