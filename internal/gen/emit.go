package gen

import (
	"iter"

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

		var file *j.File

		if p := req.ImplOutputPath(); true {
			var ok bool
			file, ok = files[p]
			if !ok {
				file = j.NewFilePathName(p, "internal")
				files[p] = file
			}
		}

		r, err := render.NewEnum(req, req.Enum)
		if err != nil {
			return err
		}

		ss, err := r.Gen()
		if err != nil {
			return err
		}
		for _, s := range ss {
			file.Add(s)
		}
	}

	return nil
}
