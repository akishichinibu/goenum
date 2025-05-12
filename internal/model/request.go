package model

import (
	"path/filepath"
	"strings"
)

type GenRequest struct {
	Unit *GenUnit
	Enum *Enum
}

func (r *GenRequest) ImplOutputPath() string {
	ext := filepath.Ext(r.Unit.Path)
	if ext == "" {
		return ""
	}
	return strings.TrimSuffix(r.Unit.Path, ext) + "_gen" + ext
}
