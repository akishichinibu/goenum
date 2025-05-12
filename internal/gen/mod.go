package gen

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

const modFileName = "go.mod"

func loadModFile(workdir string) (*modfile.File, error) {
	data, err := os.ReadFile(path.Join(workdir, modFileName))
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod: %w", err)
	}
	return modfile.Parse(modFileName, data, nil)
}

func loadPackages(workdir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: 0 |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Dir: workdir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	return pkgs, nil
}
