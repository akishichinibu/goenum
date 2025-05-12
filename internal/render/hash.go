package render

import (
	"sort"
	"strings"

	"crypto/md5"
	"encoding/hex"

	"github.com/akishichinibu/goenum/internal/model"
)

type HashString string

func (h HashString) Hash() string {
	hash := md5.Sum([]byte(h))
	return hex.EncodeToString(hash[:])
}

func hashEnum(req *model.GenRequest, e *model.Enum) (HashString, error) {
	var vs []*model.Variant
	vs = append(vs, e.Variants...)

	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Name < vs[j].Name
	})

	var tags []string
	for _, v := range vs {
		vr, err := newVariantRenderer(req, v, "")
		if err != nil {
			return "", err
		}
		tags = append(tags, vr.fingerPrint.Hash())
	}

	return HashString(strings.Join(tags, ";")), nil
}

func hashVariant(req *model.GenRequest, variant *model.Variant) (HashString, error) {
	var parts []string
	for _, param := range variant.Params {
		tr, err := model.NewType(param.Type)
		if err != nil {
			return "", err
		}
		part := param.Name + "," + tr.FingerPrint.Hash()
		parts = append(parts, part)
	}

	s := HashString(strings.Join(parts, ";"))
	return s, nil
}

func hashType(unit *model.GenUnit, t model.Type) (HashString, error) {
	path, name, err := resolveGoTypeName(unit, t)
	if err != nil {
		return "", err
	}
	if path == nil {
		var temp = "#"
		path = &temp
	}

	s := HashString(*path + ":" + name)
	return s, nil
}
