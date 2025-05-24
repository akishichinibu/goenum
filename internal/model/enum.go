package model

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

const EnumDefinitionPrefix = "_E_"

type Enum struct {
	Spec      *ast.TypeSpec
	Variants  []*Variant
	ValueType Type
}

func (e *Enum) AddVariant(variant *Variant) error {
	e.Variants = append(e.Variants, variant)

	return nil
}

func parseDeclAsEnum(unit *GenUnit, decl *ast.GenDecl, spec *ast.TypeSpec, iface *ast.InterfaceType) (e *Enum, err error) {
	if !strings.HasPrefix(spec.Name.Name, EnumDefinitionPrefix) {
		return nil, &InvalidEnumError{
			unit:   unit,
			decl:   decl,
			reason: fmt.Sprintf("goenum should be declared with `%s`, got %s", EnumDefinitionPrefix, spec.Name.Name),
		}
	}

	e = &Enum{
		Spec:      spec,
		Variants:  make([]*Variant, 0),
		ValueType: nil,
	}

	for _, method := range iface.Methods.List {
		for _, name := range method.Names {
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     spec,
					funcType: funcType,
					name:     name.Name,
					reason:   "invalid method type",
				}
			}

			variants, err := newVariant(unit, e, method, funcType)
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

	switch {
	case allValueTypeNil(e.Variants):
	case allValueTypeNotNil(e.Variants):
		for _, variant := range e.Variants {
			if len(variant.Spec.Params.List) != 0 {
				return nil, &InvalidVariantError{
					unit:     unit,
					decl:     spec,
					funcType: variant.Spec,
					name:     variant.Name,
					reason:   "value enum should not have parameters",
				}
			}

			e.ValueType = e.Variants[0].ValueType
		}
	default:
		return nil, &InvalidEnumError{
			unit:   unit,
			decl:   decl,
			reason: "enum should have all variants with same value type",
		}
	}

	return e, nil
}

func allValueTypeNil(s []*Variant) bool {
	for _, v := range s {
		if v.ValueType != nil {
			return false
		}
	}

	return true
}

func allValueTypeNotNil(s []*Variant) bool {
	for _, v := range s {
		if v.ValueType == nil {
			return false
		}
	}

	return true
}

func NewEnum(unit *GenUnit, decl *ast.GenDecl) (ms []*Enum, err error) {
	if decl.Tok != token.TYPE {
		return nil, &InvalidEnumError{
			decl:   decl,
			unit:   unit,
			reason: fmt.Sprintf("goenum should be declared with `%s`, got %s", token.TYPE, decl.Tok),
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

		meta, err := parseDeclAsEnum(unit, decl, typeSpec, iface)
		if err != nil {
			return nil, err
		}

		ms = append(ms, meta)
	}

	return ms, nil
}
