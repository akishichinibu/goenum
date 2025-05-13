package render

import (
	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

type EnumExportRenderer struct {
	req    *model.GenRequest
	Enum   *model.Enum
	naming *naming
}

func NewEnumExportRenderer(req *model.GenRequest, em *model.Enum) (Renderer, error) {
	return &EnumExportRenderer{
		req:    req,
		Enum:   em,
		naming: newNaming(em),
	}, nil
}

func (e *EnumExportRenderer) Render(emit Emitter) error {
	emit(j.Type().Id(e.naming.Interface).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.Interface))
	emit(j.Line())

	emit(j.Type().Id(e.naming.Builder).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.Builder))
	emit(j.Line())

	for _, variant := range e.Enum.Variants {
		name := e.naming.VariantInterfaceName(variant)
		emit(j.Type().Id(name).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), name))
		emit(j.Line())
	}

	emit(j.Line())
	emit(j.Var().Id(e.naming.BuilderSingleton).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.BuilderSingleton))
	emit(j.Line())

	return nil
}
