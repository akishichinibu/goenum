package gen

import (
	"go/ast"
	"iter"
	"log/slog"
	"strings"

	"github.com/akishichinibu/goenum/internal/model"
)

func scanUnits(workdir string) iter.Seq2[*model.GenUnit, error] {
	return func(yield func(*model.GenUnit, error) bool) {
		modFile, err := loadModFile(workdir)
		if err != nil {
			if !yield(nil, err) {
				return
			}
		}

		logger.Info("mod file loaded: ", modFile.Go.Version)

		packages, err := loadPackages(workdir)
		if err != nil {
			if !yield(nil, err) {
				return
			}
		}

		logger.Info("packages loaded: ", len(packages))

		for _, pkg := range packages {
			slog.Info("loading units for package: ", slog.String("pkgPath", pkg.PkgPath))
			for _, syntax := range pkg.Syntax {
				if syntax == nil {
					continue
				}

				pos := pkg.Fset.Position(syntax.Pos())

				u := &model.GenUnit{
					Path:       pos.Filename,
					ImportPath: pkg.PkgPath,
					Package:    pkg,
					Node:       syntax,
				}

				logger.Debug("found unit", slog.String("path", u.Path), slog.String("importPath", u.ImportPath))

				if !yield(u, nil) {
					return
				}
			}
		}
	}
}

func isGoEnum(decl *ast.GenDecl) bool {
	for _, comment := range decl.Doc.List {
		if strings.Contains(comment.Text, "goenum") {
			return true
		}
	}
	return false
}

func scanDecl(units iter.Seq2[*model.GenUnit, error]) iter.Seq2[*model.GenRequest, error] {
	return func(yield func(*model.GenRequest, error) bool) {
		for unit, err := range units {
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			for _, decl := range unit.Node.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				logger.Debug("found declaration", slog.String("pos", unit.Path))

				if !isGoEnum(genDecl) {
					continue
				}

				enums, err := model.NewEnum(genDecl)
				if err != nil {
					if !yield(nil, err) {
						return
					}
				}

				for _, enum := range enums {
					req := &model.GenRequest{
						Unit: unit,
						Enum: enum,
					}
					if !yield(req, nil) {
						return
					}
				}
			}
		}
	}
}
