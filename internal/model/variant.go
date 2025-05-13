package model

import (
	"go/ast"
	"strings"
)

type Variant struct {
	Enum *Enum
	Name string
	Spec *ast.FuncType

	Params    []*Param
	Value     string
	ValueType Type
}

func (v *Variant) AddParam(param *Param) error {
	v.Params = append(v.Params, param)

	return nil
}

func newVariant(unit *GenUnit, enum *Enum, field *ast.Field, funcType *ast.FuncType) (v []*Variant, err error) {
	for _, name := range field.Names {
		name := name.Name

		variant := &Variant{
			Enum: enum,
			Name: name,
			Spec: funcType,

			Params:    make([]*Param, 0),
			ValueType: nil,
			Value:     "",
		}

		if funcType.Results != nil && len(funcType.Results.List) > 0 {
			if len(funcType.Params.List) > 0 {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     enum.Spec,
					funcType: funcType,
					name:     name,
					reason:   "value enum should not have parameters",
				}
			}

			if len(funcType.Results.List) != 1 {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     enum.Spec,
					funcType: funcType,
					name:     name,
					reason:   "value enum should have exactly one return value",
				}
			}

			ret := funcType.Results.List[0]
			if len(ret.Names) > 1 {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     enum.Spec,
					funcType: funcType,
					name:     name,
					reason:   "value enum should have exactly one return value",
				}
			}

			name := ret.Names[0].Name
			if !strings.HasPrefix(name, "_") {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     enum.Spec,
					funcType: funcType,
					name:     name,
					reason:   "return value should start with `_`",
				}
			}

			value := name[1:]
			variant.Value = value

			variant.ValueType, err = NewType(ret.Type)
			if err != nil {
				return nil, err
			}
		}

		v = append(v, variant)
	}

	return v, nil
}
