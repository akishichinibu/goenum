package gen

import (
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
)

func walk(root string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		err := filepath.WalkDir(
			root,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if !yield(path, nil) {
					return filepath.SkipAll
				}
				return nil
			},
		)
		if err != nil {
			switch err {
			case filepath.SkipDir:
				return
			default:
				yield("", fmt.Errorf("error occurred while walking the directory: %w", err))
			}
		}
	}
}
