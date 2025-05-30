package render

import (
	"fmt"
	"go/types"

	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

type TypeRenderer struct {
	req         *model.GenRequest
	Type        model.Type
	FingerPrint HashString
}

func (t *TypeRenderer) Render(emit Emitter) error {
	path, name, err := resolveGoTypeName(t.req.Unit, t.Type)
	if err != nil {
		return err
	}

	if path != nil {
		emit(j.Qual(*path, name))
	} else {
		emit(j.Id(name))
	}

	return nil
}

func NewTypeRenderer(req *model.GenRequest, t model.Type) (*TypeRenderer, error) {
	fingerPrint, err := hashType(req.Unit, t)
	if err != nil {
		return nil, err
	}

	return &TypeRenderer{
		req:         req,
		Type:        t,
		FingerPrint: fingerPrint,
	}, nil
}

func resolveGoTypeName(unit *model.GenUnit, t model.Type) (path *string, name string, err error) {
	switch tt := t.(type) {
	case *model.TypeDirect:
		t := unit.Package.TypesInfo.TypeOf(tt.Indent)
		if t == nil {
			return nil, "", fmt.Errorf("cannot find type of %s", tt.Indent)
		}

		named, ok := t.(*types.Named)
		if !ok {
			return nil, t.String(), nil
		}

		obj := named.Obj()
		pkgPath := ""

		if obj.Pkg() != nil {
			pkgPath = obj.Pkg().Path()
		}

		return &pkgPath, obj.Name(), nil
	default:
		return nil, "", fmt.Errorf("unsupported type %T", t)
	}
}
