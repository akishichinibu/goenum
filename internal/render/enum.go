package render

import (
	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

type Enum struct {
	req             *model.GenRequest
	Enum            *model.Enum
	naming          *naming
	enumFingerPrint HashString
}

func NewEnum(req *model.GenRequest, em *model.Enum) (Renderer, error) {
	enumFingerPrint, err := hashEnum(req, em)
	if err != nil {
		return nil, err
	}

	return &Enum{
		Enum: em,

		req:             req,
		naming:          newNaming(em),
		enumFingerPrint: enumFingerPrint,
	}, nil
}

func (e *Enum) Render(emit Emitter) error {
	return multiRender(
		emit,
		e.renderInterface,
		e.renderVariant,
		e.renderBuilderInterface,
		e.renderBuilderImpl,
	)
}

// ExportInterface is the interface that is generated to present this enum.
func (e *Enum) renderInterface(emit Emitter) error {
	emit(
		j.Type().Id(e.naming.Interface).Interface(
			j.Id(e.naming.HashMethodName(e.enumFingerPrint)).Params(),
			j.Id(e.naming.EqualMethodName).Params(j.Id("other").Id(e.naming.Interface)).Bool(),
			j.Id(e.naming.MatchMethodName).ParamsFunc(func(p *j.Group) {
				for _, variant := range e.Enum.Variants {
					p.Id(variant.Name).Func().Params(j.Id(e.naming.VariantInterfaceName(variant)))
				}
			}),
		),
	)

	return nil
}

func (e *Enum) renderVariant(emit Emitter) error {
	for _, variant := range e.Enum.Variants {
		ev, err := newVariantRenderer(e.req, variant, e.enumFingerPrint)
		if err != nil {
			return err
		}

		if err := ev.Render(emit); err != nil {
			return err
		}
	}

	return nil
}

func (e *Enum) renderBuilderInterface(emit Emitter) error {
	emit(j.Commentf("the builder interface for %s", e.naming.Interface))
	emit(j.Line())

	methods := make([]*j.Statement, 0)

	for _, variant := range e.Enum.Variants {
		vr, err := newVariantRenderer(e.req, variant, e.enumFingerPrint)
		if err != nil {
			return err
		}

		m, err := vr.genBuilderInterfaceMethod()
		if err != nil {
			return err
		}

		methods = append(methods, m...)
	}

	emit(
		j.Type().Id(e.naming.Builder).InterfaceFunc(func(g *j.Group) {
			for _, m := range methods {
				g.Add(m)
				g.Line()
			}
		}),
	)

	return nil
}

func (e *Enum) renderBuilderImpl(emit Emitter) error {
	emit(j.Commentf("the builder for %s", e.naming.Builder))
	emit(j.Line())

	emit(j.Type().Id(e.naming.BuilderImpl).Struct())

	for _, variant := range e.Enum.Variants {
		if s, err := e.renderBuilderMethod(emit, variant); err != nil {
			return err
		} else {
			emit(s...)
		}

		emit(j.Line())
	}

	// Generate the singleton instance
	emit(j.Commentf("the singleton instance for %s", e.naming.BuilderSingleton))
	emit(j.Line())
	emit(j.Var().Id(e.naming.BuilderSingleton).Id(e.naming.Builder).Op("=").Op("&").Id(e.naming.BuilderImpl).Values())

	return nil
}

func (e *Enum) renderBuilderMethod(emit Emitter, variant *model.Variant) (ss []*j.Statement, err error) {
	vr, err := newVariantRenderer(e.req, variant, e.enumFingerPrint)
	if err != nil {
		return nil, err
	}

	pts, err := vr.resolveParamsWithType()
	if err != nil {
		return nil, err
	}

	emit(j.Commentf("the builder method for %s", e.naming.VariantBuilderName(variant)))
	emit(j.Line())

	emit(j.Func().
		Params(j.Id(e.naming.BuilderImpl)).
		Id(variant.Name).
		ParamsFunc(func(p *j.Group) {
			for _, param := range pts {
				p.Id(e.naming.ParamsPublicFieldName(param.Param)).Add(toCode(param.Statements)...)
			}
		}).
		Params(j.Id(e.naming.Interface)).
		BlockFunc(func(g *j.Group) {
			g.Return(j.Op("&").Id(e.naming.VariantImplName(variant)).ValuesFunc(func(v *j.Group) {
				for _, param := range variant.Params {
					v.Id(e.naming.ParamsPublicFieldName(param))
				}
			}))
		}),
	)

	return ss, nil
}
