package gen

import (
	"iter"
	"os"
	"path/filepath"

	"github.com/akishichinibu/goenum/internal/model"
	"github.com/akishichinibu/goenum/internal/render"
	j "github.com/dave/jennifer/jen"
)

func emitEnum(requests iter.Seq2[*model.GenRequest, error]) error {
	files := make(map[string]*j.File)

	for req, err := range requests {
		if err != nil {
			return err
		}

		var implFile *j.File

		if p := req.Unit.InternalImplFilePath(); true {
			var ok bool
			implFile, ok = files[p]
			if !ok {
				implFile = j.NewFilePathName(req.Unit.InternalImplImportPath(), req.Unit.GenPackageName())
				files[p] = implFile
			}
		}

		r, err := render.NewEnum(req, req.Enum)
		if err != nil {
			return err
		}

		ss, err := r.Render()
		if err != nil {
			return err
		}

		for _, s := range ss {
			implFile.Add(s)
		}

		var exportFile *j.File

		if p := req.Unit.ExportFilePath(); true {
			var ok bool
			exportFile, ok = files[p]
			if !ok {
				exportFile = j.NewFilePathName(req.Unit.ImportPath, req.Unit.PackageName())
				files[p] = exportFile
			}
		}

		r, err = render.NewEnumExportRenderer(req, req.Enum)
		if err != nil {
			return err
		}

		ss, err = r.Render()
		if err != nil {
			return err
		}

		for _, s := range ss {
			exportFile.Add(s)
		}
	}

	for path, file := range files {
		if err := safeSave(file, path); err != nil {
			return err
		}
	}

	return nil
}

func safeSave(f *j.File, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	if err := f.Save(path); err != nil {
		return err
	}
	return nil
}
