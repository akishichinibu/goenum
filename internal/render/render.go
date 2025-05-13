package render

import j "github.com/dave/jennifer/jen"

type Emitter func(s ...*j.Statement)

type Renderable func(Emitter) error

func multiRender(emit Emitter, es ...Renderable) error {
	for _, e := range es {
		if err := e(emit); err != nil {
			return err
		}
	}

	return nil
}

type Renderer interface {
	Render(emit Emitter) error
}

func toCode(ss []*j.Statement) []j.Code {
	codes := make([]j.Code, len(ss))
	for i, s := range ss {
		codes[i] = s
	}

	return codes
}
