package render

import (
	"github.com/akishichinibu/godantic/internal/gen/naming"
	typerenderer "github.com/akishichinibu/godantic/internal/gen/renderer/type"
	"github.com/akishichinibu/godantic/internal/model/enum"
	"github.com/akishichinibu/godantic/internal/utilities"
	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

type ParamWithType struct {
	Param      *enum.EnumParamMeta
	Type       model.Type
	Statements []*j.Statement
}

type VariantRenderer struct {
	req     *model.GenRequest
	Variant *enum.EnumVariantMeta

	naming          *naming.EnumNaming
	enumFingerPrint utilities.HashString
	fingerPrint     utilities.HashString
}

func newVariantRenderer(req *model.GenRequest, variant *enum.EnumVariantMeta, enumFingerPrint utilities.HashString) (*VariantRenderer, error) {
	fingerPrint, err := hashVariant(req, variant)
	if err != nil {
		return nil, err
	}
	return &VariantRenderer{
		req:             req,
		Variant:         variant,
		naming:          naming.NewEnumNaming(variant.Enum),
		enumFingerPrint: enumFingerPrint,
		fingerPrint:     fingerPrint,
	}, nil
}

func (e *VariantRenderer) resolveParamsWithType() ([]*ParamWithType, error) {
	params := make([]*ParamWithType, 0)
	for _, param := range e.Variant.Params {
		tr, err := typerenderer.New(e.req, param.Type)
		if err != nil {
			return nil, err
		}
		typeStatements, err := tr.Gen()
		if err != nil {
			return nil, err
		}
		params = append(params, &ParamWithType{
			Param:      param,
			Type:       param.Type,
			Statements: typeStatements,
		})
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
					p.Id(e.naming.ParamsPrivateFieldName(pt.Param)).Add(utilities.ToCode(pt.Statements)...)
				}
			}).
			Id(e.naming.Interface),
	)
	return ss, nil
}

func (e *VariantRenderer) genInterfaceMethod() ([]*j.Statement, error) {
	pts, err := e.resolveParamsWithType()
	if err != nil {
		return nil, err
	}
	ss := make([]*j.Statement, 0)
	for _, pt := range pts {
		ss = append(ss, j.Id(e.naming.ParamsGetterName(pt.Param)).Params().Add(utilities.ToCode(pt.Statements)...))
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
				g.Id(e.naming.ParamsPrivateMemberInVariant(p.Param)).Add(utilities.ToCode(p.Statements)...)
			}
		}),
	)
	return nil
}

func (e *VariantRenderer) genFieldGetters(emit Emitter) error {
	emit(j.Commentf("the getters for variant %s", e.Variant.Name))
	emit(j.Line())

	pts, err := e.resolveParamsWithType()
	if err != nil {
		return err
	}

	for _, pt := range pts {
		emit(
			j.Func().Params(j.Id("v").Op("*").Id(e.naming.VariantImplName(e.Variant))).
				Id(e.naming.ParamsGetterName(pt.Param)).
				Params().
				Params(j.Id(e.naming.ParamsReturnValueName(pt.Param)).Add(utilities.ToCode(pt.Statements)...)).
				BlockFunc(func(g *j.Group) {
					g.Return(j.Id("v").Dot(e.naming.ParamsPrivateFieldName(pt.Param)))
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

func (e *VariantRenderer) Gen() ([]*j.Statement, error) {
	return ChainRender(
		e.genInterface,
		e.genImplStruct,
		e.genFieldGetters,
		e.genHashTagImplement,
		e.genEqualImpl,
		e.genMatch,
	)
}
