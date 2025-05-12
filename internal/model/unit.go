package model

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type GenUnit struct {
	Path       string
	ImportPath string
	Package    *packages.Package
	Node       *ast.File
}

func (u *GenUnit) PathDir() string {
	return filepath.Dir(u.Path)
}

func (u *GenUnit) PackageName() string {
	return u.Node.Name.Name
}

func (u *GenUnit) GenPackageName() string {
	return fmt.Sprintf("%sgen", u.PackageName())
}

func (u *GenUnit) GodanticImplImportPath() string {
	return filepath.Join(u.ImportPath, "internal", "gen")
}

func (u *GenUnit) FileNameBase() string {
	fn := filepath.Base(u.Path)
	ext := filepath.Ext(fn)
	base := fn[:len(fn)-len(ext)]
	return base
}
