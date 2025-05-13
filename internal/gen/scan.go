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

		Logger.Info("mod file loaded", slog.String("goVersion", modFile.Go.Version))

		packages, err := loadPackages(workdir)
		if err != nil {
			if !yield(nil, err) {
				return
			}
		}

		Logger.Info("packages loaded", slog.Int("count", len(packages)))

		for _, pkg := range packages {
			slog.Info("loading units for package", slog.String("pkgPath", pkg.PkgPath))

			for _, syntax := range pkg.Syntax {
				pos := pkg.Fset.Position(syntax.Pos())

				if strings.HasSuffix(pos.Filename, "_test.go") || strings.HasSuffix(pos.Filename, ".gen.go") {
					continue
				}

				u := &model.GenUnit{
					Path:       pos.Filename,
					ImportPath: pkg.PkgPath,
					Package:    pkg,
					Node:       syntax,
				}

				Logger.Info("found unit", slog.String("path", u.Path), slog.String("importPath", u.ImportPath))

				if !yield(u, nil) {
					return
				}
			}
		}
	}
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

				Logger.Info("found declaration", slog.String("pos", unit.Path))

				if !isGoEnum(genDecl) {
					Logger.Info("not a godantic:enum", slog.String("pos", unit.Path))

					continue
				}

				Logger.Info("found godantic:enum", slog.String("pos", unit.Path))

				enums, err := model.NewEnum(unit, genDecl)
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

func isGoEnum(decl *ast.GenDecl) bool {
	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			if strings.Contains(comment.Text, "godantic:enum") {
				return true
			}
		}
	}

	return false
}
