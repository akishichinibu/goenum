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

func NewEnumExportRenderer(req *model.GenRequest, em *model.Enum) *EnumExportRenderer {
	return &EnumExportRenderer{
		req:    req,
		Enum:   em,
		naming: newNaming(em),
	}
}

func (e *EnumExportRenderer) Gen() (ss []*j.Statement, err error) {
	ss = append(ss, j.Type().Id(e.naming.Interface).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.Interface))
	ss = append(ss, j.Line())

	ss = append(ss, j.Type().Id(e.naming.Builder).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.Builder))
	ss = append(ss, j.Line())

	for _, variant := range e.Enum.Variants {
		name := e.naming.VariantInterfaceName(variant)
		ss = append(ss, j.Type().Id(name).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), name))
		ss = append(ss, j.Line())
	}

	ss = append(ss, j.Line())
	ss = append(ss, j.Var().Id(e.naming.BuilderSingleton).Op("=").Qual(e.req.Unit.GodanticImplImportPath(), e.naming.BuilderSingleton))
	ss = append(ss, j.Line())
	return ss, nil
}
