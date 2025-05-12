package render

import (
	"fmt"
	"go/types"

	"github.com/akishichinibu/goenum/internal/model"
	j "github.com/dave/jennifer/jen"
)

func resolveGoTypeName(unit *model.GenUnit, t model.Type) (path *string, name string, err error) {
	switch tt := t.(type) {
	case *model.TypeDirect:
		if tt.Indent == nil {
			return nil, "", fmt.Errorf("indent is nil")
		}

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
	case *model.ArrayType:
		path, name, err = resolveGoTypeName(unit, tt.ElemType)
		if err != nil {
			return nil, "", err
		}
		return path, "[]" + name, nil
	default:
		return nil, "", fmt.Errorf("unsupported type %T", t)
	}
}

type TypeRenderer struct {
	req         *model.GenRequest
	Type        model.Type
	FingerPrint HashString
}

func (t *TypeRenderer) Gen() (ss []*j.Statement, err error) {
	path, name, err := resolveGoTypeName(t.req.Unit, t.Type)
	if err != nil {
		return nil, err
	}
	if path != nil {
		ss = append(ss, j.Qual(*path, name))
	} else {
		ss = append(ss, j.Id(name))
	}
	return ss, nil
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
