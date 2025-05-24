package render

import (
	"fmt"
	"strconv"

	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

type ParamWithType struct {
	Param      *model.Param
	Type       model.Type
	Statements []*j.Statement
}

type VariantRenderer struct {
	req     *model.GenRequest
	Variant *model.Variant

	naming          *naming
	enumFingerPrint HashString
	fingerPrint     HashString
}

func newVariantRenderer(req *model.GenRequest, variant *model.Variant, enumFingerPrint HashString) (*VariantRenderer, error) {
	fingerPrint, err := hashVariant(req, variant)
	if err != nil {
		return nil, err
	}

	return &VariantRenderer{
		req:             req,
		Variant:         variant,
		naming:          newNaming(variant.Enum),
		enumFingerPrint: enumFingerPrint,
		fingerPrint:     fingerPrint,
	}, nil
}

func (e *VariantRenderer) Render(emit Emitter) error {
	return multiRender(
		emit,
		e.genInterface,
		e.genImplStruct,
		e.genFieldGetters,
		e.genHashTagImplement,
		e.genEqualImpl,
		e.genMatch,
		e.genValue,
	)
}

func (e *VariantRenderer) resolveParamsWithType() (params []*ParamWithType, error error) {
	for _, param := range e.Variant.Params {
		pt := &ParamWithType{
			Param:      param,
			Type:       param.Type,
			Statements: make([]*j.Statement, 0),
		}

		tr, err := NewTypeRenderer(e.req, param.Type)
		if err != nil {
			return nil, err
		}

		if err := tr.Render(func(s ...*j.Statement) {
			pt.Statements = append(pt.Statements, s...)
		}); err != nil {
			return nil, err
		}

		params = append(params, pt)
	}

	return params, nil
}

func (e *VariantRenderer) genBuilderInterfaceMethod() (ss []*j.Statement, err error) {
	pts, err := e.resolveParamsWithType()
	if err != nil {
		return nil, err
	}

	ss = append(
		ss,
		j.Id(e.Variant.Name).
			ParamsFunc(func(p *j.Group) {
				for _, pt := range pts {
					p.Id(e.naming.ParamsPrivateFieldName(pt.Param)).Add(toCode(pt.Statements)...)
				}
			}).
			Id(e.naming.Interface),
	)

	return ss, nil
}

func (e *VariantRenderer) genInterfaceMethod() (ss []*j.Statement, err error) {
	pts, err := e.resolveParamsWithType()
	if err != nil {
		return nil, err
	}

	for _, pt := range pts {
		ss = append(ss, j.Id(e.naming.ParamsGetterName(pt.Param)).Params().Add(toCode(pt.Statements)...))
	}

	return ss, nil
}

func (e *VariantRenderer) genInterface(emit Emitter) error {
	methods, err := e.genInterfaceMethod()
	if err != nil {
		return err
	}

	emit(j.Commentf("the interface for variant %s", e.Variant.Name))
	emit(j.Line())
	emit(j.Type().Id(e.naming.VariantInterfaceName(e.Variant)).InterfaceFunc(func(g *j.Group) {
		g.Id("_enum_" + e.enumFingerPrint.Hash()).Params()

		for _, m := range methods {
			g.Add(m)
			g.Line()
		}
	}))

	return nil
}

func (e *VariantRenderer) genImplStruct(emit Emitter) error {
	pts, err := e.resolveParamsWithType()
	if err != nil {
		return err
	}

	emit(
		j.Commentf("the implementation struct for variant %s", e.Variant.Name),
		j.Line(),
		j.Type().Id(e.naming.VariantImplName(e.Variant)).StructFunc(func(g *j.Group) {
			for _, p := range pts {
				g.Id(e.naming.ParamsPrivateMemberInVariant(p.Param)).Add(toCode(p.Statements)...)
			}
		}),
	)

	return nil
}

func (e *VariantRenderer) genFieldGetters(emit Emitter) error {
	emit(j.Commentf("the getters for variant %s", e.Variant.Name))
	emit(j.Line())

	params, err := e.resolveParamsWithType()
	if err != nil {
		return err
	}

	for _, param := range params {
		emit(
			j.Func().Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
				Id(e.naming.ParamsGetterName(param.Param)).
				Params().
				Params(j.Id(e.naming.ParamsReturnValueName(param.Param)).Add(toCode(param.Statements)...)).
				BlockFunc(func(g *j.Group) {
					g.Return(j.Id("v").Dot(e.naming.ParamsPrivateFieldName(param.Param)))
				}),
			j.Line(),
		)
	}

	return nil
}

func (e *VariantRenderer) genHashTagImplement(emit Emitter) error {
	emit(
		j.Func().Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
			Id("_enum_" + e.enumFingerPrint.Hash()).
			Params().
			Block(),
	)

	return nil
}

func (e *VariantRenderer) genEqualImpl(emit Emitter) error {
	emit(
		j.Commentf("the Equal method for variant %s", e.Variant.Name),
		j.Line(),
		j.Func().
			Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
			Id(e.naming.EqualMethodName).
			Params(j.Id("other").Id(e.naming.Interface)).
			Bool().
			BlockFunc(func(g *j.Group) {
				g.If(j.Id("v").Op("==").Nil().Op("||").Id("other").Op("==").Nil()).Block(
					j.Return(j.False()),
				)
				g.List(j.Id("otherImpl"), j.Id("ok")).Op(":=").Id("other").Assert(j.Op("*").Id(e.naming.VariantImplName(e.Variant)))
				g.If(j.Op("!").Id("ok")).Block(
					j.Return(j.False()),
				)
				g.Return(j.Op("*").Id("v").Op("==").Op("*").Id("otherImpl"))
			}),
	)

	return nil
}

func (e *VariantRenderer) genMatch(emit Emitter) error {
	emit(j.Commentf("the Match method for variant %s", e.Variant.Name))
	emit(j.Line())
	emit(j.Func().Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
		Id(e.naming.MatchMethodName).
		ParamsFunc(func(p *j.Group) {
			for _, v := range e.Variant.Enum.Variants {
				p.Id(v.Name).Func().Params(j.Id(e.naming.VariantInterfaceName(v)))
			}
		}).
		BlockFunc(func(g *j.Group) {
			g.Id(e.Variant.Name).Call(j.Id("v"))
		}),
	)

	return nil
}

func (e *VariantRenderer) genValue(emit Emitter) error {
	if e.Variant.ValueType == nil {
		return nil
	}

	_, name, err := resolveGoTypeName(e.req.Unit, e.Variant.ValueType)
	if err != nil {
		return err
	}

	var retValue any

	var retType string

	switch name {
	case "string":
		retType = "string"
		retValue = e.Variant.Value
	case "int":
		retType = "int"

		retValue, err = strconv.ParseInt(e.Variant.Value, 10, 0)
		if err != nil {
			return fmt.Errorf("cannot parse int value %s: %w", e.Variant.Value, err)
		}

		retValue = int(retValue.(int64))
	default:
		return fmt.Errorf("unsupported value type %s", name)
	}

	emit(
		j.Commentf("the Value method for variant %s", e.Variant.Name),
		j.Line(),
		j.Func().
			Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
			Id(e.naming.ValueMethodName).
			Params().
			Id(retType).
			BlockFunc(func(g *j.Group) {
				g.Return(j.Lit(retValue))
			}),
	)

	return nil
}
